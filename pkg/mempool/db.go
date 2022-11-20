package mempool

import (
	"encoding/json"
	"strings"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

const keySeparator = ":"
const keyPrefix = "mempool"

func getDBKey(entryPoint common.Address, sender common.Address) []byte {
	return []byte(keyPrefix + keySeparator + entryPoint.String() + keySeparator + sender.String())
}

func getEntryPointAndSenderFromDBKey(key []byte) (common.Address, common.Address) {
	slc := strings.Split(string(key), keySeparator)
	ep := common.HexToAddress(slc[1])
	sender := common.HexToAddress(slc[2])

	return ep, sender
}

func getUserOpFromDBValue(value []byte) (*userop.UserOperation, error) {
	data := make(map[string]any)
	json.Unmarshal(value, &data)
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
			ep, _ := getEntryPointAndSenderFromDBKey(item.Key())

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
