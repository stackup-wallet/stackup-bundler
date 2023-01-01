package checks

import (
	"errors"

	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// ValidatePendingOps checks the pending UserOperations by the same sender and only passes if:
//
//  1. Sender doesn't have another UserOperation already present in the pool.
//  2. It replaces an existing UserOperation with same nonce and higher fee.
//  3. Sender is staked and is allowed multiple UserOperations in the pool.
func ValidatePendingOps(op *userop.UserOperation, penOps []*userop.UserOperation, gs GetStakeFunc) error {
	_, err := gs(op.Sender)
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
			if op.MaxPriorityFeePerGas.Cmp(oldOp.MaxPriorityFeePerGas) <= 0 {
				return errors.New(
					"pending ops: sender has op in mempool with same or higher priority fee",
				)
			}
		}
		// else if !dep.Staked {
		// 	return errors.New(
		// 		"pending ops: sender must be staked to have multiple ops in the mempool",
		// 	)
		// }
	}
	return nil
}
