package standalone

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"golang.org/x/sync/errgroup"
)

// SanityCheck returns a UserOpHandler that relies on a given ethClient to run through all the standard
// client checks as specified in EIP-4337. This should be the first module in the stack.
func SanityCheck(eth *ethclient.Client, maxVerificationGas *big.Int) modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		ep, err := entrypoint.NewEntrypoint(ctx.EntryPoint, eth)
		if err != nil {
			return err
		}

		g := new(errgroup.Group)
		g.Go(func() error { return checkSender(eth, ctx.UserOp) })
		g.Go(func() error { return checkVerificationGas(maxVerificationGas, ctx.UserOp) })
		g.Go(func() error { return checkPaymasterAndData(eth, ep, ctx.UserOp) })
		g.Go(func() error { return checkCallGasLimit(ctx.UserOp) })
		g.Go(func() error { return checkFeePerGas(eth, ctx.UserOp) })

		if err := g.Wait(); err != nil {
			return err
		}
		return nil
	}
}

// Simulation returns a UserOpHandler that relies on a given ethClient to run through the standard simulation
// as specified in EIP-4337. This should be done after all checks are complete.
func Simulation(eth *ethclient.Client) modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		_, err := entrypoint.SimulateValidation(eth, ctx.EntryPoint, ctx.UserOp)
		return err
	}
}

// TrackPaymasterDeposit returns a BatchHandler that tracks each paymaster in the batch and ensures it has
// enough deposit to pay for all the UserOps that use it.
func TrackPaymasterDeposit(eth *ethclient.Client) modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		ep, err := entrypoint.NewEntrypoint(ctx.EntryPoint, eth)
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
