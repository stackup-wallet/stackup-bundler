// Package checks implements modules for running an array of standard validations for both the Client and
// Bundler.
package checks

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"golang.org/x/sync/errgroup"
)

// Standalone exposes modules to perform basic Client and Bundler checks as specified in EIP-4337. It is
// intended for bundlers that are independent of an Ethereum node and hence relies on a given ethClient to
// query blockchain state.
type Standalone struct {
	rpc                *rpc.Client
	eth                *ethclient.Client
	maxVerificationGas *big.Int
	tracer             string
}

// New returns a Standalone instance with methods that can be used in Client and Bundler modules to perform
// standard checks as specified in EIP-4337.
func New(rpc *rpc.Client, maxVerificationGas *big.Int, tracer string) *Standalone {
	eth := ethclient.NewClient(rpc)
	return &Standalone{rpc, eth, maxVerificationGas, tracer}
}

// ValidateOpValues returns a UserOpHandler that runs through some first line sanity checks for new UserOps
// received by the Client. This should be one of the first modules executed by the Client.
func (s *Standalone) ValidateOpValues() modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		penOps := ctx.GetPendingOps()
		gc := getCodeWithEthClient(s.eth)
		gbf := getBaseFeeWithEthClient(s.eth)
		gs, err := getStakeWithEthClient(ctx, s.eth)
		if err != nil {
			return err
		}

		g := new(errgroup.Group)
		g.Go(func() error { return ValidateSender(ctx.UserOp, gc) })
		g.Go(func() error { return ValidateInitCode(ctx.UserOp, gs) })
		g.Go(func() error { return ValidateVerificationGas(ctx.UserOp, s.maxVerificationGas) })
		g.Go(func() error { return ValidatePaymasterAndData(ctx.UserOp, gc, gs) })
		g.Go(func() error { return ValidateCallGasLimit(ctx.UserOp) })
		g.Go(func() error { return ValidateFeePerGas(ctx.UserOp, gbf) })
		g.Go(func() error { return ValidatePendingOps(ctx.UserOp, penOps, gs) })

		if err := g.Wait(); err != nil {
			return errors.NewRPCError(errors.INVALID_FIELDS, err.Error(), err.Error())
		}
		return nil
	}
}

// SimulateOp returns a UserOpHandler that runs through simulation of new UserOps with the EntryPoint.
func (s *Standalone) SimulateOp() modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		g := new(errgroup.Group)
		g.Go(func() error {
			_, err := entrypoint.SimulateValidation(s.rpc, ctx.EntryPoint, ctx.UserOp)

			if err != nil {
				return errors.NewRPCError(errors.REJECTED_BY_EP_OR_ACCOUNT, err.Error(), err.Error())
			}
			return nil
		})
		g.Go(func() error {
			err := entrypoint.TraceSimulateValidation(
				s.rpc,
				ctx.EntryPoint,
				ctx.UserOp,
				ctx.ChainID,
				s.tracer,
				entrypoint.EntityStakes{
					ctx.UserOp.GetFactory():   ctx.GetDepositInfo(ctx.UserOp.GetFactory()),
					ctx.UserOp.Sender:         ctx.GetDepositInfo(ctx.UserOp.Sender),
					ctx.UserOp.GetPaymaster(): ctx.GetDepositInfo(ctx.UserOp.GetPaymaster()),
				},
			)

			if err != nil {
				return errors.NewRPCError(errors.BANNED_OPCODE, err.Error(), err.Error())
			}
			return nil
		})

		return g.Wait()
	}
}

// PaymasterDeposit returns a BatchHandler that tracks each paymaster in the batch and ensures it has enough
// deposit to pay for all the UserOps that use it.
func (s *Standalone) PaymasterDeposit() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		ep, err := entrypoint.NewEntrypoint(ctx.EntryPoint, s.eth)
		if err != nil {
			return err
		}

		deps := make(map[common.Address]*big.Int)
		for i, op := range ctx.Batch {
			pm := op.GetPaymaster()
			if pm == common.HexToAddress("0x") {
				continue
			}

			if _, ok := deps[pm]; !ok {
				dep, err := ep.GetDepositInfo(nil, pm)
				if err != nil {
					return err
				}

				deps[pm] = dep.Deposit
			}

			deps[pm] = big.NewInt(0).Sub(deps[pm], op.GetMaxPrefund())
			if deps[pm].Cmp(common.Big0) < 0 {
				ctx.MarkOpIndexForRemoval(i)
			}
		}

		return nil
	}
}
