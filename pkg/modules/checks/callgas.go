package checks

import (
	"fmt"

	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// ValidateCallGasLimit checks the callGasLimit is at least the cost of a CALL with non-zero value.
func ValidateCallGasLimit(op *userop.UserOperation) error {
	cg := gas.NewDefaultOverhead().NonZeroValueCall()
	if op.CallGasLimit.Cmp(cg) < 0 {
		return fmt.Errorf("callGasLimit: below expected gas of %s", cg.String())
	}

	return nil
}
