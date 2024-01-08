package entities

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/dbutils"
)

type addressCounter map[common.Address]int

type status int64

const (
	ok status = iota
	throttled
	banned
)

var (
	emaHours       = 24
	opsCountPrefix = dbutils.JoinValues("entity", "opsCount")
)

func getOpsCountKey(entity common.Address) []byte {
	return []byte(dbutils.JoinValues(opsCountPrefix, entity.String()))
}

func getOpsCountValue(opsSeen int, opsIncluded int) []byte {
	return []byte(
		dbutils.JoinValues(strconv.Itoa(opsSeen), strconv.Itoa(opsIncluded), fmt.Sprint(time.Now().Unix())),
	)
}

func applyExpWeights(txn *badger.Txn, key []byte, value []byte) (opsSeen int, opsIncluded int, err error) {
	counts := dbutils.SplitValues(string(value))
	opsSeen, err = strconv.Atoi(counts[0])
	if err != nil {
		return 0, 0, err
	}
	opsIncluded, err = strconv.Atoi(counts[1])
	if err != nil {
		return 0, 0, err
	}
	lastUpdated, err := strconv.ParseInt(counts[2], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	dur := time.Since(time.Unix(lastUpdated, 0))
	for i := int(dur.Hours()); i > 0; i-- {
		if opsSeen < 24 && opsIncluded < 24 {
			break
		}

		opsSeen -= opsSeen / emaHours
		opsIncluded -= opsIncluded / emaHours
	}

	e := badger.NewEntry(key, getOpsCountValue(opsSeen, opsIncluded))
	err = txn.SetEntry(e)

	return opsSeen, opsIncluded, err
}

func getOpsCountByEntity(
	txn *badger.Txn,
	entity common.Address,
) (opsSeen int, opsIncluded int, err error) {
	key := getOpsCountKey(entity)
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

	return applyExpWeights(txn, key, value)
}

func incrementOpsSeenByEntity(txn *badger.Txn, entity common.Address) error {
	opsSeen, opsIncluded, err := getOpsCountByEntity(txn, entity)
	if err != nil {
		return err
	}

	e := badger.NewEntry(getOpsCountKey(entity), getOpsCountValue(opsSeen+1, opsIncluded))
	return txn.SetEntry(e)
}

func incrementOpsIncludedByEntity(txn *badger.Txn, count addressCounter) error {
	for entity, n := range count {
		opsSeen, opsIncluded, err := getOpsCountByEntity(txn, entity)
		if err != nil {
			return err
		}

		e := badger.NewEntry(
			getOpsCountKey(entity),
			getOpsCountValue(opsSeen, opsIncluded+n),
		)
		if err := txn.SetEntry(e); err != nil {
			return err
		}
	}

	return nil
}

func getStatus(txn *badger.Txn, entity common.Address, repConst *ReputationConstants) (status, error) {
	opsSeen, opsIncluded, err := getOpsCountByEntity(txn, entity)
	if err != nil {
		return ok, err
	}
	if opsSeen == 0 {
		return ok, nil
	}

	minExpectedIncluded := opsSeen / repConst.MinInclusionRateDenominator
	if minExpectedIncluded <= opsIncluded+repConst.ThrottlingSlack {
		return ok, nil
	} else if minExpectedIncluded <= opsIncluded+repConst.BanSlack {
		return throttled, nil
	} else {
		return banned, nil
	}
}

func overrideEntity(txn *badger.Txn, entry *ReputationOverride) error {
	return txn.SetEntry(
		badger.NewEntry(getOpsCountKey(entry.Address), getOpsCountValue(entry.OpsSeen, entry.OpsIncluded)),
	)
}
