package checks

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Checks the callGasLimit is at least the cost of a CALL with non-zero value.
func checkCallGasLimit(op *userop.UserOperation) error {
	cg := gas.NewDefaultOverhead().NonZeroValueCall()
	if op.CallGasLimit.Cmp(cg) < 0 {
		return fmt.Errorf("callGasLimit: below expected gas of %s", cg.String())
	}

	return nil
}

// The maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client
// is willing to accept. At the minimum, they are sufficiently high to be included with the current
// block.basefee.
func checkFeePerGas(eth *ethclient.Client, op *userop.UserOperation) error {
	tip, err := eth.SuggestGasTipCap(context.Background())
	if err != nil {
		return err
	}

	if op.MaxPriorityFeePerGas.Cmp(tip) < 0 {
		return fmt.Errorf("maxPriorityFeePerGas: below expected wei of %s", tip.String())
	}
	if op.MaxFeePerGas.Cmp(op.MaxPriorityFeePerGas) < 0 {
		return fmt.Errorf("maxFeePerGas: must be equal to or greater than maxPriorityFeePerGas")
	}

	return nil
}
