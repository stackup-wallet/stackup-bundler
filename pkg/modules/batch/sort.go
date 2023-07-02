package batch

import (
	"sort"

	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
)

// SortByNonce returns a BatchHandlerFunc that ensures ops with same sender is ordered by ascending nonce
// regardless of gas price.
func SortByNonce() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		sort.SliceStable(ctx.Batch, func(i, j int) bool {
			return ctx.Batch[i].Sender == ctx.Batch[j].Sender &&
				ctx.Batch[i].Nonce.Cmp(ctx.Batch[j].Nonce) == -1
		})

		return nil
	}
}
