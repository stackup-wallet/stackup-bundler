package paymaster

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/dbutils"
)

type addressCounter map[string]int

type status int64

const (
	ok status = iota
	throttled
	banned
)

const minInclusionRateDenominator = 10
const throttlingSlack = 10
const banSlack = 50
const emaHours = 24

var (
	opsCountPrefix = dbutils.JoinValues("paymaster", "opsCount")
)

func getOpsCountKey(paymaster common.Address) []byte {
	return []byte(dbutils.JoinValues(opsCountPrefix, paymaster.String()))
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

func getOpsCountByPaymaster(
	txn *badger.Txn,
	paymaster common.Address,
) (opsSeen int, opsIncluded int, err error) {
	key := getOpsCountKey(paymaster)
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

func incrementOpsSeenByPaymaster(txn *badger.Txn, paymaster common.Address) error {
	opsSeen, opsIncluded, err := getOpsCountByPaymaster(txn, paymaster)
	if err != nil {
		return err
	}

	e := badger.NewEntry(getOpsCountKey(paymaster), getOpsCountValue(opsSeen+1, opsIncluded))
	return txn.SetEntry(e)
}

func incrementOpsIncludedByPaymasters(
	txn *badger.Txn,
	count addressCounter,
	paymasters ...common.Address,
) error {
	for _, paymaster := range paymasters {
		opsSeen, opsIncluded, err := getOpsCountByPaymaster(txn, paymaster)
		if err != nil {
			return err
		}

		e := badger.NewEntry(
			getOpsCountKey(paymaster),
			getOpsCountValue(opsSeen, opsIncluded+count[paymaster.String()]),
		)
		if err := txn.SetEntry(e); err != nil {
			return err
		}
	}

	return nil
}

func getStatus(txn *badger.Txn, paymaster common.Address) (status, error) {
	opsSeen, opsIncluded, err := getOpsCountByPaymaster(txn, paymaster)
	if err != nil {
		return ok, err
	}
	if opsSeen == 0 {
		return ok, nil
	}

	minExpectedIncluded := opsSeen / minInclusionRateDenominator
	if minExpectedIncluded <= opsIncluded+throttlingSlack {
		return ok, nil
	} else if minExpectedIncluded <= opsIncluded+banSlack {
		return throttled, nil
	} else {
		return banned, nil
	}
}
