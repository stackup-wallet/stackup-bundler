package checks

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

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
