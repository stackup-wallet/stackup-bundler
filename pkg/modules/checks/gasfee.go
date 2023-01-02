package checks

import (
	"fmt"

	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// ValidateFeePerGas checks the maxFeePerGas is sufficiently high to be included with the current
// block.basefee.
func ValidateFeePerGas(op *userop.UserOperation, gbf GetBaseFeeFunc) error {
	bf, err := gbf()
	if err != nil {
		return err
	}

	if op.MaxFeePerGas.Cmp(bf) < 0 {
		return fmt.Errorf("maxFeePerGas: must be equal to or greater than current block.basefee")
	}

	return nil
}
