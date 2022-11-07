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

func getOpFunc(q *userOpQueues) GetOp {
	return func(entryPoint common.Address, sender common.Address) (*userop.UserOperation, error) {
		op := q.GetOp(entryPoint, sender)
		return op, nil
	}
}

func addOpFunc(db *badger.DB, q *userOpQueues) AddOp {
	return func(entryPoint common.Address, op *userop.UserOperation) error {
		data, err := op.MarshalJSON()
		if err != nil {
			return err
		}

		err = db.Update(func(txn *badger.Txn) error {
			return txn.Set(getDBKey(entryPoint, op.Sender), data)
		})
		if err != nil {
			return err
		}

		q.AddOp(entryPoint, op)
		return nil
	}
}

func bundleOpsFunc(db *badger.DB, q *userOpQueues) BundleOps {
	return func(entryPoint common.Address) ([]*userop.UserOperation, error) {
		return q.Next(entryPoint), nil
	}
}

func removeOpsFunc(db *badger.DB, q *userOpQueues) RemoveOps {
	return func(entryPoint common.Address, senders ...common.Address) error {
		err := db.Update(func(txn *badger.Txn) error {
			for _, s := range senders {
				err := txn.Delete(getDBKey(entryPoint, s))
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		q.RemoveOps(entryPoint, senders)
		return nil
	}
}

// NewBadgerDBWrapper creates an instance of a mempool that uses badgerDB
// to persist and load UserOperations from disk incase of a reset.
func NewBadgerDBWrapper(db *badger.DB) (*Interface, error) {
	q := newUserOpQueue()
	err := loadFromDisk(db, q)
	if err != nil {
		return nil, err
	}

	return &Interface{
		AddOp:     addOpFunc(db, q),
		GetOp:     getOpFunc(q),
		BundleOps: bundleOpsFunc(db, q),
		RemoveOps: removeOpsFunc(db, q),
	}, nil
}
