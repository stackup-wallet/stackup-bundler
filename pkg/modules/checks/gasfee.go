package checks

import (
	"fmt"

	"github.com/stackup-wallet/stackup-bundler/pkg/modules/gasprice"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// ValidateFeePerGas checks the maxFeePerGas is sufficiently high to be included with the current
// block.basefee. Alternatively, if basefee is not supported, then check that maxPriorityFeePerGas is equal to
// maxFeePerGas as a fallback.
func ValidateFeePerGas(op *userop.UserOperation, gbf gasprice.GetBaseFeeFunc) error {
	bf, err := gbf()
	if err != nil {
		return err
	}

	if bf == nil {
		if op.MaxPriorityFeePerGas.Cmp(op.MaxFeePerGas) != 0 {
			return fmt.Errorf("legacy fee mode: maxPriorityFeePerGas must equal maxFeePerGas")
		}

		return nil
	}

	if op.MaxPriorityFeePerGas.Cmp(op.MaxFeePerGas) == 1 {
		return fmt.Errorf("maxFeePerGas: must be equal to or greater than maxPriorityFeePerGas")
	}

	if op.MaxFeePerGas.Cmp(bf) < 0 {
		return fmt.Errorf("maxFeePerGas: must be equal to or greater than current block.basefee")
	}

	return nil
}
