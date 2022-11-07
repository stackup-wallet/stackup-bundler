package relay

import (
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

const banThreshold = 3
const timeWindow = 7 * 24 * time.Hour

const separator = ":"
const keyPrefix = "module" + separator + "relay"
const opsCountPrefix = keyPrefix + separator + "opsCount"
const requestIDPrefix = keyPrefix + separator + "requestID"

func getOpsCountKey(clientID string) []byte {
	return []byte(opsCountPrefix + separator + clientID)
}

func getRequestIDKey(requestID string) []byte {
	return []byte(requestIDPrefix + separator + requestID)
}

func getRequestIDsFromOps(ep common.Address, chainID *big.Int, ops ...*userop.UserOperation) []string {
	rids := []string{}
	for _, op := range ops {
		rids = append(rids, op.GetRequestID(ep, chainID).String())
	}

	return rids
}

func getOpsCountByKey(txn *badger.Txn, key []byte) (opsSeen int, opsIncluded int, err error) {
	item, err := txn.Get(key)
	if err != nil && err == badger.ErrKeyNotFound {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	var value []byte
	err = item.Value(func(val []byte) error {
		value = append([]byte{}, val...)
		return nil
	})
	if err != nil {
		return 0, 0, err
	}

	counts := strings.Split(string(value), separator)
	opsSeen, err = strconv.Atoi(counts[0])
	if err != nil {
		return 0, 0, err
	}
	opsIncluded, err = strconv.Atoi(counts[1])
	if err != nil {
		return 0, 0, err
	}

	return opsSeen, opsIncluded, nil
}

func getOpsCount(c *gin.Context, txn *badger.Txn, clientID string) (opsSeen int, opsIncluded int, err error) {
	return getOpsCountByKey(txn, getOpsCountKey(clientID))
}

func incrementOpsSeen(c *gin.Context, txn *badger.Txn, clientID string) error {
	opsSeen, opsIncluded, err := getOpsCount(c, txn, clientID)
	if err != nil {
		return err
	}

	val := strconv.Itoa(opsSeen+1) + separator + strconv.Itoa(opsIncluded)
	e := badger.NewEntry(getOpsCountKey(clientID), []byte(val)).WithTTL(timeWindow)
	return txn.SetEntry(e)
}

func incrementOpsIncludedByRequestIDs(txn *badger.Txn, requestIDs ...string) error {
	for _, rid := range requestIDs {
		item, err := txn.Get(getRequestIDKey(rid))
		if err != nil {
			return err
		}

		var clientID []byte
		err = item.Value(func(val []byte) error {
			clientID = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}

		opsCountKey := getOpsCountKey(string(clientID))
		opsSeen, opsIncluded, err := getOpsCountByKey(txn, opsCountKey)
		if err != nil {
			return err
		}

		var opsCountValue string
		if opsSeen == 0 && opsIncluded == 0 {
			// Op has been in the mempool longer than timeWindow
			opsCountValue = strconv.Itoa(opsSeen+1) + separator + strconv.Itoa(opsIncluded+1)
		} else {
			opsCountValue = strconv.Itoa(opsSeen) + separator + strconv.Itoa(opsIncluded+1)
		}

		e := badger.NewEntry(opsCountKey, []byte(opsCountValue)).WithTTL(timeWindow)
		if err := txn.SetEntry(e); err != nil {
			return err
		}
	}

	return nil
}

func mapRequestIDToClient(c *gin.Context, txn *badger.Txn, requestID string, clientID string) error {
	return txn.Set(getRequestIDKey(requestID), []byte(clientID))
}

func removeRequestIDEntries(txn *badger.Txn, requestIDs ...string) error {
	for _, rid := range requestIDs {
		err := txn.Delete(getRequestIDKey(rid))
		if err != nil {
			return err
		}
	}

	return nil
}
