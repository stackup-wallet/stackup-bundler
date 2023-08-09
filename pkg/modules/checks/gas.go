package checks

import (
	"fmt"
	"math/big"

	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// ValidateGasAvailable checks that the max available gas is less than the batch gas limit.
func ValidateGasAvailable(op *userop.UserOperation, maxBatchGasLimit *big.Int) error {
	// This calculation ensures that we are only checking the gas used for execution. In rollups, the PVG also
	// includes the L1 callData cost. If the L1 gas component spikes, it can cause the PVG value of legit ops
	// to be greater than the maxBatchGasLimit. For non-rollups, the results would be the same as just calling
	// op.GetMaxGasAvailable().
	static, err := gas.NewDefaultOverhead().CalcPreVerificationGas(op)
	if err != nil {
		return err
	}
	mgl := big.NewInt(0).Sub(op.GetMaxGasAvailable(), op.PreVerificationGas)
	mga := big.NewInt(0).Add(mgl, static)

	if mga.Cmp(maxBatchGasLimit) > 0 {
		return fmt.Errorf("gasLimit: exceeds maxBatchGasLimit of %s", maxBatchGasLimit.String())
	}

	return nil
}
