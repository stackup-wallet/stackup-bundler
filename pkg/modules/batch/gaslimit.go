package batch

import (
	"math/big"

	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// MaintainGasLimit returns a BatchHandlerFunc that ensures the max gas used from the entire batch does not
// exceed the allowed threshold.
func MaintainGasLimit(maxBatchGasLimit *big.Int) modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		bat := []*userop.UserOperation{}
		sum := big.NewInt(0)
		for _, op := range ctx.Batch {
			sum = big.NewInt(0).Add(sum, op.GetMaxGasAvailable())
			if sum.Cmp(maxBatchGasLimit) >= 0 {
				break
			}
			bat = append(bat, op)
		}
		ctx.Batch = bat

		return nil
	}
}
