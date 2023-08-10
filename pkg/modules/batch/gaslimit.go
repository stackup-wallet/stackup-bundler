package batch

import (
	"math/big"

	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// MaintainGasLimit returns a BatchHandlerFunc that ensures the max gas used from the entire batch does not
// exceed the allowed threshold.
func MaintainGasLimit(maxBatchGasLimit *big.Int) modules.BatchHandlerFunc {
	// See comment in pkg/modules/checks/gas.go
	staticOv := gas.NewDefaultOverhead()

	return func(ctx *modules.BatchHandlerCtx) error {
		bat := []*userop.UserOperation{}
		sum := big.NewInt(0)
		for _, op := range ctx.Batch {
			static, err := staticOv.CalcPreVerificationGas(op)
			if err != nil {
				return err
			}
			mgl := big.NewInt(0).Sub(op.GetMaxGasAvailable(), op.PreVerificationGas)
			mga := big.NewInt(0).Add(mgl, static)

			sum = big.NewInt(0).Add(sum, mga)
			if sum.Cmp(maxBatchGasLimit) >= 0 {
				break
			}
			bat = append(bat, op)
		}
		ctx.Batch = bat

		return nil
	}
}
