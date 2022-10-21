package mempool

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type PendingTransaction struct {
	UserOp     *userop.UserOperation
	EntryPoint common.Address
}

type Add func(sender string, op *userop.UserOperation, entryPoint common.Address) (bool, error)
type Get func(sender string) (*PendingTransaction, error)
type Batch func(size int) ([]*PendingTransaction, error)

type ClientInterface struct {
	Add Add
	Get Get
}

type BundlerInterface struct {
	Batch Batch
}
