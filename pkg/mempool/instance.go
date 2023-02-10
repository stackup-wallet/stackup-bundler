// Package mempool provides a local representation of all the UserOperations that are known to the bundler
// which have passed all Client checks and pending action by the Bundler.
package mempool

import (
	badger "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type watch struct {
	key string
	sig chan bool
}

// Mempool provides read and write access to a pool of pending UserOperations which have passed all Client
// checks.
type Mempool struct {
	db      *badger.DB
	queue   *userOpQueues
	watches []watch
}

// New creates an instance of a mempool that uses an embedded DB to persist and load UserOperations from disk
// incase of a reset.
func New(db *badger.DB) (*Mempool, error) {
	queue := newUserOpQueue()
	watches := []watch{}
	err := loadFromDisk(db, queue)
	if err != nil {
		return nil, err
	}

	return &Mempool{db, queue, watches}, nil
}

func (m *Mempool) pushSignals() {
	for _, watch := range m.watches {
		watch.sig <- true
	}
}

// GetOps returns all the UserOperations associated with an EntryPoint and Sender address.
func (m *Mempool) GetOps(entryPoint common.Address, sender common.Address) ([]*userop.UserOperation, error) {
	ops := m.queue.GetOps(entryPoint, sender)
	return ops, nil
}

// AddOp adds a UserOperation to the mempool or replace an existing one with the same EntryPoint, Sender, and
// Nonce values.
func (m *Mempool) AddOp(entryPoint common.Address, op *userop.UserOperation) error {
	data, err := op.MarshalJSON()
	if err != nil {
		return err
	}

	err = m.db.Update(func(txn *badger.Txn) error {
		return txn.Set(getUniqueKey(entryPoint, op.Sender, op.Nonce), data)
	})
	if err != nil {
		return err
	}

	m.queue.AddOp(entryPoint, op)
	m.pushSignals()
	return nil
}

// BundleOps builds a bundle of UserOperations from the mempool to be sent to the EntryPoint.
func (m *Mempool) BundleOps(entryPoint common.Address) ([]*userop.UserOperation, error) {
	return m.queue.Next(entryPoint), nil
}

// RemoveOps removes a list of UserOperations from the mempool by EntryPoint, Sender, and Nonce values.
func (m *Mempool) RemoveOps(entryPoint common.Address, ops ...*userop.UserOperation) error {
	err := m.db.Update(func(txn *badger.Txn) error {
		for _, op := range ops {
			err := txn.Delete(getUniqueKey(entryPoint, op.Sender, op.Nonce))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	m.queue.RemoveOps(entryPoint, ops...)
	return nil
}

// OnAdd allows entities to register a channel that will receive a signal every time a UserOperation is added.
// It returns a function to stop receiving.
func (m *Mempool) OnAdd(sig chan bool) func() {
	key := uuid.New().String()
	w := watch{key, sig}
	m.watches = append(m.watches, w)

	// initial push if queue is non-empty
	if m.queue.opCount > 0 {
		w.sig <- true
	}

	return func() {
		f := []watch{}
		for _, w := range m.watches {
			if w.key != key {
				f = append(f, w)
			}
		}
		m.watches = f
	}
}

// Dump will return a list of UserOperations from the mempool by EntryPoint in the order it arrived.
func (m *Mempool) Dump(entryPoint common.Address) ([]*userop.UserOperation, error) {
	return m.queue.All(entryPoint), nil
}

// Clear will clear the entire embedded db and reset it to a clean state.
func (m *Mempool) Clear() error {
	if err := m.db.DropAll(); err != nil {
		return err
	}
	m.queue = newUserOpQueue()

	return nil
}
