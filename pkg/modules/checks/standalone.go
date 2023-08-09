// Package checks implements modules for running an array of standard validations for both the Client and
// Bundler.
package checks

import (
	"math/big"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/simulation"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/gasprice"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"golang.org/x/sync/errgroup"
)

// Standalone exposes modules to perform basic Client and Bundler checks as specified in EIP-4337. It is
// intended for bundlers that are independent of an Ethereum node and hence relies on a given ethClient to
// query blockchain state.
type Standalone struct {
	db                      *badger.DB
	rpc                     *rpc.Client
	eth                     *ethclient.Client
	ov                      *gas.Overhead
	maxVerificationGas      *big.Int
	maxBatchGasLimit        *big.Int
	maxOpsForUnstakedSender int
}

// New returns a Standalone instance with methods that can be used in Client and Bundler modules to perform
// standard checks as specified in EIP-4337.
func New(
	db *badger.DB,
	rpc *rpc.Client,
	ov *gas.Overhead,
	maxVerificationGas *big.Int,
	maxBatchGasLimit *big.Int,
	maxOpsForUnstakedSender int,
) *Standalone {
	eth := ethclient.NewClient(rpc)
	return &Standalone{db, rpc, eth, ov, maxVerificationGas, maxBatchGasLimit, maxOpsForUnstakedSender}
}

// ValidateOpValues returns a UserOpHandler that runs through some first line sanity checks for new UserOps
// received by the Client. This should be one of the first modules executed by the Client.
func (s *Standalone) ValidateOpValues() modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		penOps := ctx.GetPendingOps()
		gc := getCodeWithEthClient(s.eth)
		gbf := gasprice.GetBaseFeeWithEthClient(s.eth)
		gs, err := getStakeWithEthClient(ctx, s.eth)
		if err != nil {
			return err
		}

		g := new(errgroup.Group)
		g.Go(func() error { return ValidateSender(ctx.UserOp, gc) })
		g.Go(func() error { return ValidateInitCode(ctx.UserOp, gs) })
		g.Go(func() error { return ValidateVerificationGas(ctx.UserOp, s.ov, s.maxVerificationGas) })
		g.Go(func() error { return ValidatePaymasterAndData(ctx.UserOp, gc, gs) })
		g.Go(func() error { return ValidateCallGasLimit(ctx.UserOp, s.ov) })
		g.Go(func() error { return ValidateFeePerGas(ctx.UserOp, gbf) })
		g.Go(func() error { return ValidatePendingOps(ctx.UserOp, penOps, s.maxOpsForUnstakedSender, gs) })
		g.Go(func() error { return ValidateGasAvailable(ctx.UserOp, s.maxBatchGasLimit) })

		if err := g.Wait(); err != nil {
			return errors.NewRPCError(errors.INVALID_FIELDS, err.Error(), err.Error())
		}
		return nil
	}
}

// SimulateOp returns a UserOpHandler that runs through simulation of new UserOps with the EntryPoint.
func (s *Standalone) SimulateOp() modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		gc := getCodeWithEthClient(s.eth)
		g := new(errgroup.Group)
		g.Go(func() error {
			sim, err := simulation.SimulateValidation(s.rpc, ctx.EntryPoint, ctx.UserOp)

			if err != nil {
				return errors.NewRPCError(errors.REJECTED_BY_EP_OR_ACCOUNT, err.Error(), err.Error())
			}
			if sim.ReturnInfo.SigFailed {
				return errors.NewRPCError(
					errors.INVALID_SIGNATURE,
					"Invalid UserOp signature or paymaster signature",
					nil,
				)
			}
			if sim.ReturnInfo.ValidUntil.Cmp(common.Big0) != 0 &&
				time.Now().Unix() >= sim.ReturnInfo.ValidUntil.Int64()-30 {
				return errors.NewRPCError(
					errors.SHORT_DEADLINE,
					"expires too soon",
					nil,
				)
			}
			return nil
		})
		g.Go(func() error {
			ic, err := simulation.TraceSimulateValidation(
				s.rpc,
				ctx.EntryPoint,
				ctx.UserOp,
				ctx.ChainID,
				simulation.EntityStakes{
					ctx.UserOp.GetFactory():   ctx.GetDepositInfo(ctx.UserOp.GetFactory()),
					ctx.UserOp.Sender:         ctx.GetDepositInfo(ctx.UserOp.Sender),
					ctx.UserOp.GetPaymaster(): ctx.GetDepositInfo(ctx.UserOp.GetPaymaster()),
				},
			)
			if err != nil {
				return errors.NewRPCError(errors.BANNED_OPCODE, err.Error(), err.Error())
			}

			ch, err := getCodeHashes(ic, gc)
			if err != nil {
				return errors.NewRPCError(errors.BANNED_OPCODE, err.Error(), err.Error())
			}
			return saveCodeHashes(s.db, ctx.UserOp.GetUserOpHash(ctx.EntryPoint, ctx.ChainID), ch)
		})

		return g.Wait()
	}
}

// CodeHashes returns a BatchHandler that verifies the code for any interacted contracts has not changed since
// the first simulation.
func (s *Standalone) CodeHashes() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		gc := getCodeWithEthClient(s.eth)

		end := len(ctx.Batch) - 1
		for i := end; i >= 0; i-- {
			op := ctx.Batch[i]
			chs, err := getSavedCodeHashes(s.db, op.GetUserOpHash(ctx.EntryPoint, ctx.ChainID))
			if err != nil {
				return err
			}

			changed, err := hasCodeHashChanges(chs, gc)
			if err != nil {
				return err
			}
			if changed {
				ctx.MarkOpIndexForRemoval(i)
			}
		}
		return nil
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

// Clean returns a BatchHandler that clears the DB of data that is no longer required. This should be one of
// the last modules executed by the Bundler.
func (s *Standalone) Clean() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		all := append([]*userop.UserOperation{}, ctx.Batch...)
		all = append(all, ctx.PendingRemoval...)
		hashes := []common.Hash{}
		for _, op := range all {
			hashes = append(hashes, op.GetUserOpHash(ctx.EntryPoint, ctx.ChainID))
		}

		return removeSavedCodeHashes(s.db, hashes...)
	}
}
