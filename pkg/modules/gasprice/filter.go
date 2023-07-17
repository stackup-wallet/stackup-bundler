package gasprice

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// FilterUnderpriced returns a BatchHandlerFunc that will filter out all the userOps that are below either the
// dynamic or legacy GasPrice set in the context.
func FilterUnderpriced() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		b := []*userop.UserOperation{}
		for _, op := range ctx.Batch {
			if ctx.BaseFee != nil && ctx.BaseFee.Cmp(common.Big0) != 0 && ctx.Tip != nil {
				gp := big.NewInt(0).Add(ctx.BaseFee, ctx.Tip)
				if op.GetDynamicGasPrice(ctx.BaseFee).Cmp(gp) >= 0 {
					b = append(b, op)
				}
			} else if ctx.GasPrice != nil {
				if op.MaxFeePerGas.Cmp(ctx.GasPrice) >= 0 {
					b = append(b, op)
				}
			}
		}

		ctx.Batch = b
		return nil
	}
}
