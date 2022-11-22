package modules

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// NewBatchHandlerContext creates a new BatchHandlerCtx using a copy of the given batch.
func NewBatchHandlerContext(batch []*userop.UserOperation, entryPoint common.Address, chainID *big.Int) *BatchHandlerCtx {
	var copy []*userop.UserOperation
	copy = append(copy, batch...)

	return &BatchHandlerCtx{
		Batch:          copy,
		PendingRemoval: []*userop.UserOperation{},
		EntryPoint:     entryPoint,
		ChainID:        chainID,
		Data:           make(map[string]any),
	}
}

// NewUserOpHandlerContext creates a new UserOpHandlerCtx using a given op.
func NewUserOpHandlerContext(op *userop.UserOperation, entryPoint common.Address, chainID *big.Int) *UserOpHandlerCtx {
	return &UserOpHandlerCtx{
		UserOp:     op,
		EntryPoint: entryPoint,
		ChainID:    chainID,
	}
}
