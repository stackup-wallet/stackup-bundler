package mempool

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// GetOp checks if a UserOperation is in the mempool and returns it.
type GetOp func(entryPoint common.Address, sender common.Address) (*userop.UserOperation, error)

// AddOp adds a UserOperation to the mempool.
type AddOp func(entryPoint common.Address, op *userop.UserOperation) error

// BundleOps builds a bundle of ops from the mempool to be sent to the EntryPoint.
type BundleOps func(entryPoint common.Address) ([]*userop.UserOperation, error)

// RemoveOps removes a list of UserOperations from the mempool by sender.
type RemoveOps func(entryPoint common.Address, senders []common.Address) error

type Interface struct {
	GetOp
	AddOp
	BundleOps
	RemoveOps
}
