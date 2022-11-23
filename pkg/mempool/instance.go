// Package mempool provides a local representation of all the UserOperations that are known to the bundler
// which have passed all Client checks and pending action by the Bundler.
package mempool

import (
	badger "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Mempool provides read and write access to a pool of pending UserOperations which have passed all Client
// checks.
type Mempool struct {
	db    *badger.DB
	queue *userOpQueues
}

// New creates an instance of a mempool that uses an embedded DB to persist and load UserOperations from disk
// incase of a reset.
func New(db *badger.DB) (*Mempool, error) {
	queue := newUserOpQueue()
	err := loadFromDisk(db, queue)
	if err != nil {
		return nil, err
	}

	return &Mempool{db, queue}, nil
}

// GetOp checks if a UserOperation is in the mempool and returns it.
func (m *Mempool) GetOp(entryPoint common.Address, sender common.Address) (*userop.UserOperation, error) {
	op := m.queue.GetOp(entryPoint, sender)
	return op, nil
}

// AddOp adds a UserOperation to the mempool.
func (m *Mempool) AddOp(entryPoint common.Address, op *userop.UserOperation) error {
	data, err := op.MarshalJSON()
	if err != nil {
		return err
	}

	err = m.db.Update(func(txn *badger.Txn) error {
		return txn.Set(getDBKey(entryPoint, op.Sender), data)
	})
	if err != nil {
		return err
	}

	m.queue.AddOp(entryPoint, op)
	return nil
}

// BundleOps builds a bundle of ops from the mempool to be sent to the EntryPoint.
func (m *Mempool) BundleOps(entryPoint common.Address) ([]*userop.UserOperation, error) {
	return m.queue.Next(entryPoint), nil
}

// RemoveOps removes a list of UserOperations from the mempool by sender address.
func (m *Mempool) RemoveOps(entryPoint common.Address, senders ...common.Address) error {
	err := m.db.Update(func(txn *badger.Txn) error {
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

	m.queue.RemoveOps(entryPoint, senders)
	return nil
}
