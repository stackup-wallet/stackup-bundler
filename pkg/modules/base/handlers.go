package base

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// StandaloneClient returns a UserOpHandler that relies on a given ethClient to run through all the standard
// client checks as specified in EIP-4337. This should be the first module in the stack.
func StandaloneClient(eth *ethclient.Client, maxVerificationGas *big.Int) modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		ep, err := entrypoint.NewEntrypoint(ctx.EntryPoint, eth)
		if err != nil {
			return err
		}

		// Sanity checks
		if err := checkSender(eth, ctx.UserOp); err != nil {
			return err
		}
		if err := checkVerificationGas(maxVerificationGas, ctx.UserOp); err != nil {
			return err
		}
		if err := checkPaymasterAndData(eth, ctx.UserOp, ep); err != nil {
			return err
		}
		if err := checkCallGasLimit(eth, ctx.UserOp); err != nil {
			return err
		}
		if err := checkFeePerGas(eth, ctx.UserOp); err != nil {
			return err
		}

		// Op simulation
		if _, err := entrypoint.SimulateValidation(ep, entrypoint.UserOperation(*ctx.UserOp)); err != nil {
			return err
		}

		return nil
	}
}

// StandaloneBundler returns a BatchHandler that relies on a given ethClient to run through all the standard
// bundler checks as specified in EIP-4337. This should be the first module in the stack.
func StandaloneBundler(eth *ethclient.Client) modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		var filter []*userop.UserOperation
		filter = append(filter, ctx.Batch...)

		filter = filterSender(filter)
		filter = filterPaymaster(filter)

		ctx.Batch = filter
		return nil
	}
}
