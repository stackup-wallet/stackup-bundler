package bundler

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/noop"
)

// New initializes a new ERC-4337 bundler which can be extended with modules for validating batches and
// excluding UserOperations that should not be sent to the EntryPoint.
func New(mempool *mempool.Interface, chainID *big.Int, supportedEntryPoints []common.Address) *Bundler {
	return &Bundler{
		mempool:              mempool,
		chainID:              chainID,
		supportedEntryPoints: supportedEntryPoints,
		batchHandler:         noop.BatchHandler,
		errorHandler:         noop.ErrorHandler,
	}
}
