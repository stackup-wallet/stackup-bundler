package checks

import (
	"fmt"
	"math/big"

	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

const minPriceBump = 10

// calcNewThresholds returns new threshold values where newFee = oldFee  * (100 + minPriceBump) / 100.
func calcNewThresholds(cap *big.Int, tip *big.Int) (newCap *big.Int, newTip *big.Int) {
	a := big.NewInt(100 + minPriceBump)
	aFeeCap := big.NewInt(0).Mul(a, cap)
	aTip := big.NewInt(0).Mul(a, tip)

	b := big.NewInt(100)
	newCap = aFeeCap.Div(aFeeCap, b)
	newTip = aTip.Div(aTip, b)

	return newCap, newTip
}

// ValidatePendingOps checks the pending UserOperations by the same sender and only passes if:
//
//  1. Sender doesn't have another UserOperation already present in the pool.
//  2. It replaces an existing UserOperation with same nonce and higher fee.
//  3. Sender is staked and is allowed uncapped UserOperations in the pool.
func ValidatePendingOps(
	op *userop.UserOperation,
	penOps []*userop.UserOperation,
	maxOpsForUnstakedSender int,
	gs GetStakeFunc,
) error {
	dep, err := gs(op.Sender)
	if err != nil {
		return err
	}

	if len(penOps) > 0 {
		var oldOp *userop.UserOperation
		for _, penOp := range penOps {
			if op.Nonce.Cmp(penOp.Nonce) == 0 {
				oldOp = penOp
			}
		}

		if oldOp != nil {
			newMf, newMpf := calcNewThresholds(oldOp.MaxFeePerGas, oldOp.MaxPriorityFeePerGas)

			if op.MaxFeePerGas.Cmp(newMf) < 0 || op.MaxPriorityFeePerGas.Cmp(newMpf) < 0 {
				return fmt.Errorf(
					"pending ops: replacement op must increase maxFeePerGas and MaxPriorityFeePerGas by >= %d%%",
					minPriceBump,
				)
			}
		} else if !dep.Staked && len(penOps) >= maxOpsForUnstakedSender {
			return fmt.Errorf(
				"pending ops: sender must be staked to have more than %d ops in the mempool",
				maxOpsForUnstakedSender,
			)
		}
	}
	return nil
}
