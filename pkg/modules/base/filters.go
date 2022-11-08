package base

import "github.com/stackup-wallet/stackup-bundler/pkg/userop"

// Exclude UserOps that access any sender address created by another UserOp on the same batch (via CREATE2 factory).
func filterSender(batch []*userop.UserOperation) []*userop.UserOperation {
	// TODO: Add implementation
	return batch
}

// For each paymaster used in the batch, keep track of the balance while adding UserOps.
// Ensure that it has sufficient deposit to pay for all the UserOps that use it.
func filterPaymaster(batch []*userop.UserOperation) []*userop.UserOperation {
	// TODO: Add implementation
	return batch
}
