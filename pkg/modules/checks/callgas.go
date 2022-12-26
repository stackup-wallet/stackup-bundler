package checks

import (
	"fmt"

	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func ValidateCallGasLimit(op *userop.UserOperation) error {
	cg := gas.NewDefaultOverhead().NonZeroValueCall()
	if op.CallGasLimit.Cmp(cg) < 0 {
		return fmt.Errorf("callGasLimit: below expected gas of %s", cg.String())
	}

	return nil
}
