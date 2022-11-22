package checks

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"golang.org/x/sync/errgroup"
)

// Standalone exposes modules to perform basic Client and Bundler checks as specified in EIP-4337. It is
// intended for bundlers that are independent of an Ethereum node and hence relies on a given ethClient to
// query blockchain state.
type Standalone struct {
	eth                *ethclient.Client
	maxVerificationGas *big.Int
}

// New returns a Standalone instance with methods that can be used as Client and Bundler modules to perform
// standard checks as specified in EIP-4337.
func New(eth *ethclient.Client, maxVerificationGas *big.Int) *Standalone {
	return &Standalone{eth, maxVerificationGas}
}

// ValidateOpValues returns a UserOpHandler that runs through some first line sanity checks for new UserOps
// going into the Client. This should be one of the first modules in the Client stack.
func (s *Standalone) ValidateOpValues() modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		ep, err := entrypoint.NewEntrypoint(ctx.EntryPoint, s.eth)
		if err != nil {
			return err
		}

		g := new(errgroup.Group)
		g.Go(func() error { return checkSender(s.eth, ctx.UserOp) })
		g.Go(func() error { return checkVerificationGas(s.maxVerificationGas, ctx.UserOp) })
		g.Go(func() error { return checkPaymasterAndData(s.eth, ep, ctx.UserOp) })
		g.Go(func() error { return checkCallGasLimit(ctx.UserOp) })
		g.Go(func() error { return checkFeePerGas(s.eth, ctx.UserOp) })

		if err := g.Wait(); err != nil {
			return err
		}
		return nil
	}
}

// SimulateOp returns a UserOpHandler that runs through simulation of new UserOps with the EntryPoint. This
// should be done after all validations are complete.
func (s *Standalone) SimulateOp() modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		_, err := entrypoint.SimulateValidation(s.eth, ctx.EntryPoint, ctx.UserOp)
		return err
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
