package mempool

import (
	"encoding/json"
	"math/big"
	"strings"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

const keySeparator = ":"
const keyPrefix = "mempool"

func getUniqueKey(entryPoint common.Address, sender common.Address, nonce *big.Int) []byte {
	return []byte(
		keyPrefix + keySeparator + entryPoint.String() + keySeparator + sender.String() + keySeparator + nonce.String(),
	)
}

func getEntryPointFromDBKey(key []byte) common.Address {
	slc := strings.Split(string(key), keySeparator)
	return common.HexToAddress(slc[1])
}

func getUserOpFromDBValue(value []byte) (*userop.UserOperation, error) {
	data := make(map[string]any)
	if err := json.Unmarshal(value, &data); err != nil {
		return nil, err
	}

	op, err := userop.New(data)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func loadFromDisk(db *badger.DB, q *userOpQueues) error {
	return db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		prefix := []byte(keyPrefix)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			ep := getEntryPointFromDBKey(item.Key())

			err := item.Value(func(v []byte) error {
				op, err := getUserOpFromDBValue(v)
				if err != nil {
					return err
				}

				q.AddOp(ep, op)
				return nil
			})

			if err != nil {
				return err
			}
		}

		return nil
	})
}
