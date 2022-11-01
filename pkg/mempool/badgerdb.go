package mempool

import (
	"encoding/json"
	"strings"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

const keySeparator = ":"

func getDBKey(entryPoint common.Address, sender common.Address) []byte {
	return []byte(entryPoint.String() + keySeparator + sender.String())
}

func getEntryPointAndSenderFromDBKey(key []byte) (common.Address, common.Address) {
	slc := strings.Split(string(key), keySeparator)
	ep := common.HexToAddress(slc[0])
	sender := common.HexToAddress(slc[1])

	return ep, sender
}

func getUserOpFromDBValue(value []byte) (*userop.UserOperation, error) {
	data := make(map[string]interface{})
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
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
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
	return func(entryPoint common.Address, op *userop.UserOperation) (bool, error) {
		data, err := op.MarshalJSON()
		if err != nil {
			return false, nil
		}

		err = db.Update(func(txn *badger.Txn) error {
			return txn.Set(getDBKey(entryPoint, op.Sender), data)
		})
		if err != nil {
			return false, nil
		}

		return q.AddOp(entryPoint, op), nil
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
		AddOp: addOpFunc(db, q),
		GetOp: getOpFunc(q),
	}, nil
}
