package standalone

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// SanityCheck returns a UserOpHandler that relies on a given ethClient to run through all the standard
// client checks as specified in EIP-4337. This should be the first module in the stack.
func SanityCheck(eth *ethclient.Client, maxVerificationGas *big.Int) modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		ep, err := entrypoint.NewEntrypoint(ctx.EntryPoint, eth)
		if err != nil {
			return err
		}

		if err := checkSender(eth, ctx.UserOp); err != nil {
			return err
		}
		if err := checkVerificationGas(maxVerificationGas, ctx.UserOp); err != nil {
			return err
		}
		if err := checkPaymasterAndData(eth, ep, ctx.UserOp); err != nil {
			return err
		}
		if err := checkCallGasLimit(ctx.UserOp); err != nil {
			return err
		}
		if err := checkFeePerGas(eth, ctx.UserOp); err != nil {
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

// Filter returns a BatchHandler that relies on a given ethClient to run through all the standard bundler
// checks as specified in EIP-4337. This should be the first module in the stack.
func Filter(eth *ethclient.Client) modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		var filter []*userop.UserOperation
		filter = append(filter, ctx.Batch...)

		filter = filterSender(filter)
		filter = filterPaymaster(filter)

		ctx.Batch = filter
		return nil
	}
}
