package bundler

import "github.com/stackup-wallet/stackup-bundler/pkg/userop"

func adjustBatchSize(max int, batch []*userop.UserOperationV06) []*userop.UserOperationV06 {
	if len(batch) > max && max > 0 {
		return batch[:max]
	}
	return batch
}
