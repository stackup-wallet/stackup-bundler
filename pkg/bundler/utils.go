package bundler

import "github.com/stackup-wallet/stackup-bundler/pkg/userop"

func adjustBatchSize(max int, batch []*userop.UserOperation) []*userop.UserOperation {
	if len(batch) > max && max > 0 {
		return batch[:max]
	}
	return batch
}
