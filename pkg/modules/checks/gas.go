package checks

import (
	"fmt"
	"math/big"

	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// ValidateGasAvailable checks that the max available gas is less than the batch gas limit.
func ValidateGasAvailable(op *userop.UserOperation, maxBatchGasLimit *big.Int) error {
	if op.GetMaxGasAvailable().Cmp(maxBatchGasLimit) > 0 {
		return fmt.Errorf("gasLimit: exceeds maxBatchGasLimit of %s", maxBatchGasLimit.String())
	}

	return nil
}
