package modules

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/stake"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// BatchHandlerCtx is the object passed to BatchHandler functions during the Bundler's Run process. It
// also contains a Data field for adding arbitrary key-value pairs to the context. These values will be
// logged by the Bundler at the end of each run.
type BatchHandlerCtx struct {
	Batch          []*userop.UserOperation
	PendingRemoval []*userop.UserOperation
	EntryPoint     common.Address
	ChainID        *big.Int
	BaseFee        *big.Int
	Tip            *big.Int
	GasPrice       *big.Int
	Data           map[string]any
}

// NewBatchHandlerContext creates a new BatchHandlerCtx using a copy of the given batch.
func NewBatchHandlerContext(
	batch []*userop.UserOperation,
	entryPoint common.Address,
	chainID *big.Int,
	baseFee *big.Int,
	tip *big.Int,
	gasPrice *big.Int,
) *BatchHandlerCtx {
	var copy []*userop.UserOperation
	copy = append(copy, batch...)

	return &BatchHandlerCtx{
		Batch:          copy,
		PendingRemoval: []*userop.UserOperation{},
		EntryPoint:     entryPoint,
		ChainID:        chainID,
		BaseFee:        baseFee,
		Tip:            tip,
		GasPrice:       gasPrice,
		Data:           make(map[string]any),
	}
}

// MarkOpIndexForRemoval will remove the op by index from the batch and add it to the pending removal array.
// This should be used for ops that are not to be included on-chain and dropped from the mempool.
func (c *BatchHandlerCtx) MarkOpIndexForRemoval(index int) {
	batch := []*userop.UserOperation{}
	var op *userop.UserOperation
	for i, curr := range c.Batch {
		if i == index {
			op = curr
		} else {
			batch = append(batch, curr)
		}
	}
	if op == nil {
		return
	}

	c.Batch = batch
	c.PendingRemoval = append(c.PendingRemoval, op)
}

// UserOpHandlerCtx is the object passed to UserOpHandler functions during the Client's SendUserOperation
// process.
type UserOpHandlerCtx struct {
	UserOp              *userop.UserOperation
	EntryPoint          common.Address
	ChainID             *big.Int
	pendingSenderOps    []*userop.UserOperation
	pendingFactoryOps   []*userop.UserOperation
	pendingPaymasterOps []*userop.UserOperation
	senderDeposit       *entrypoint.IStakeManagerDepositInfo
	factoryDeposit      *entrypoint.IStakeManagerDepositInfo
	paymasterDeposit    *entrypoint.IStakeManagerDepositInfo
}

// NewUserOpHandlerContext creates a new UserOpHandlerCtx using a given op.
func NewUserOpHandlerContext(
	op *userop.UserOperation,
	entryPoint common.Address,
	chainID *big.Int,
	mem *mempool.Mempool,
	gs stake.GetStakeFunc,
) (*UserOpHandlerCtx, error) {
	// Fetch any pending UserOperations in the mempool by entity
	pso, err := mem.GetOps(entryPoint, op.Sender)
	if err != nil {
		return nil, err
	}
	pfo, err := mem.GetOps(entryPoint, op.GetFactory())
	if err != nil {
		return nil, err
	}
	ppo, err := mem.GetOps(entryPoint, op.GetPaymaster())
	if err != nil {
		return nil, err
	}

	// Fetch the current entrypoint deposits by entity
	sd, err := gs(entryPoint, op.Sender)
	if err != nil {
		return nil, err
	}
	fd, err := gs(entryPoint, op.GetFactory())
	if err != nil {
		return nil, err
	}
	pd, err := gs(entryPoint, op.GetPaymaster())
	if err != nil {
		return nil, err
	}

	return &UserOpHandlerCtx{
		UserOp:              op,
		EntryPoint:          entryPoint,
		ChainID:             chainID,
		pendingSenderOps:    pso,
		pendingFactoryOps:   pfo,
		pendingPaymasterOps: ppo,
		senderDeposit:       sd,
		factoryDeposit:      fd,
		paymasterDeposit:    pd,
	}, nil
}

// GetSenderDepositInfo returns the current EntryPoint deposit for the sender.
func (c *UserOpHandlerCtx) GetSenderDepositInfo() *entrypoint.IStakeManagerDepositInfo {
	return c.senderDeposit
}

// GetFactoryDepositInfo returns the current EntryPoint deposit for the factory.
func (c *UserOpHandlerCtx) GetFactoryDepositInfo() *entrypoint.IStakeManagerDepositInfo {
	return c.factoryDeposit
}

// GetPaymasterDepositInfo returns the current EntryPoint deposit for the paymaster.
func (c *UserOpHandlerCtx) GetPaymasterDepositInfo() *entrypoint.IStakeManagerDepositInfo {
	return c.paymasterDeposit
}

// GetPendingSenderOps returns all pending UserOperations in the mempool by the same sender.
func (c *UserOpHandlerCtx) GetPendingSenderOps() []*userop.UserOperation {
	return c.pendingSenderOps
}

// GetPendingFactoryOps returns all pending UserOperations in the mempool by the same factory.
func (c *UserOpHandlerCtx) GetPendingFactoryOps() []*userop.UserOperation {
	return c.pendingFactoryOps
}

// GetPendingPaymasterOps returns all pending UserOperations in the mempool by the same paymaster.
func (c *UserOpHandlerCtx) GetPendingPaymasterOps() []*userop.UserOperation {
	return c.pendingPaymasterOps
}
