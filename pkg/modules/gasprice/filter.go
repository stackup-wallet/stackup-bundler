package gasprice

import (
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// FilterUnderpriced returns a BatchHandlerFunc that will filter out all the userOps that are below either the
// BaseFee or GasPrice set in the context.
func FilterUnderpriced() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		b := []*userop.UserOperation{}
		for _, op := range ctx.Batch {
			if ctx.BaseFee != nil {
				if op.GetGasPrice(ctx.BaseFee).Cmp(ctx.BaseFee) >= 0 {
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
