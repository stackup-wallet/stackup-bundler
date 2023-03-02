package relay

import (
	"math/big"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/dbutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

const NoBanThreshold = 0
const DefaultBanThreshold = 3
const DefaultBanTimeWindow = 7 * 24 * time.Hour

var (
	keyPrefix        = dbutils.JoinValues("relay")
	opsCountPrefix   = dbutils.JoinValues(keyPrefix, "opsCount")
	userOpHashPrefix = dbutils.JoinValues(keyPrefix, "userOpHash")
)

func getOpsCountKey(clientID string) []byte {
	return []byte(dbutils.JoinValues(opsCountPrefix, clientID))
}

func getUserOpHashKey(userOpHash string) []byte {
	return []byte(dbutils.JoinValues(userOpHashPrefix, userOpHash))
}

func getUserOpHashesFromOps(ep common.Address, chainID *big.Int, ops ...*userop.UserOperation) []string {
	hashes := []string{}
	for _, op := range ops {
		hashes = append(hashes, op.GetUserOpHash(ep, chainID).String())
	}

	return hashes
}

func getOpsCountByClientID(txn *badger.Txn, clientID string) (opsSeen int, opsIncluded int, err error) {
	item, err := txn.Get(getOpsCountKey(clientID))
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

	counts := dbutils.SplitValues(string(value))
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

func incrementOpsSeenByClientID(txn *badger.Txn, clientID string, ttl time.Duration) error {
	opsSeen, opsIncluded, err := getOpsCountByClientID(txn, clientID)
	if err != nil {
		return err
	}

	val := dbutils.JoinValues(strconv.Itoa(opsSeen+1), strconv.Itoa(opsIncluded))
	e := badger.NewEntry(getOpsCountKey(clientID), []byte(val)).WithTTL(ttl)
	return txn.SetEntry(e)
}

func incrementOpsIncludedByUserOpHashes(txn *badger.Txn, ttl time.Duration, userOpHashes ...string) error {
	for _, hash := range userOpHashes {
		item, err := txn.Get(getUserOpHashKey(hash))
		if err != nil && err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil {
			return err
		}

		var value []byte
		err = item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}

		cid := string(value)
		opsSeen, opsIncluded, err := getOpsCountByClientID(txn, cid)
		if err != nil {
			return err
		}

		var opsCountValue string
		if opsSeen == 0 && opsIncluded == 0 {
			// Op has been in the mempool longer than TTL
			opsCountValue = dbutils.JoinValues(strconv.Itoa(opsSeen+1), strconv.Itoa(opsIncluded+1))
		} else {
			opsCountValue = dbutils.JoinValues(strconv.Itoa(opsSeen), strconv.Itoa(opsIncluded+1))
		}

		e := badger.NewEntry(getOpsCountKey(cid), []byte(opsCountValue)).WithTTL(ttl)
		if err := txn.SetEntry(e); err != nil {
			return err
		}
	}

	return nil
}

func mapUserOpHashToClientID(txn *badger.Txn, userOpHash string, clientID string) error {
	return txn.Set(getUserOpHashKey(userOpHash), []byte(clientID))
}

func removeUserOpHashEntries(txn *badger.Txn, userOpHashes ...string) error {
	for _, hashes := range userOpHashes {
		err := txn.Delete(getUserOpHashKey(hashes))
		if err != nil {
			return err
		}
	}

	return nil
}
