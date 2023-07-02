package gasprice

import (
	"sort"

	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
)

// SortByGasPrice returns a BatchHandlerFunc that will sort the context batch by highest GasPrice first.
func SortByGasPrice() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		if ctx.BaseFee != nil {
			sort.SliceStable(ctx.Batch, func(i, j int) bool {
				return ctx.Batch[i].GetDynamicGasPrice(ctx.BaseFee).
					Cmp(ctx.Batch[j].GetDynamicGasPrice(ctx.BaseFee)) ==
					1
			})
		} else if ctx.GasPrice != nil {
			sort.SliceStable(ctx.Batch, func(i, j int) bool {
				return ctx.Batch[i].MaxFeePerGas.Cmp(ctx.Batch[j].MaxFeePerGas) == 1
			})
		}

		return nil
	}
}
