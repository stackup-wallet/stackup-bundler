package bundler

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
)

// Bundler controls the end to end process of creating a batch of UserOperations from the mempool and sending
// it to the EntryPoint.
type Bundler struct {
	mempool              *mempool.Interface
	chainID              *big.Int
	supportedEntryPoints []common.Address
	batchHandler         modules.BatchHandlerFunc
	errorHandler         modules.ErrorHandlerFunc
}

// UseModules defines the BatchHandlers to process batches after it has gone through the standard checks.
func (i *Bundler) UseModules(handlers ...modules.BatchHandlerFunc) {
	i.batchHandler = modules.ComposeBatchHandlerFunc(handlers...)
}

// SetErrorHandlerFunc defines a method for handling errors at any point of the process.
func (i *Bundler) SetErrorHandlerFunc(handler modules.ErrorHandlerFunc) {
	i.errorHandler = handler
}

// Run starts a goroutine that will continuously process batches from the mempool.
func (i *Bundler) Run() error {
	go func(i *Bundler) {
		for {
			for _, ep := range i.supportedEntryPoints {
				batch, err := i.mempool.BundleOps(ep)
				if err != nil {
					i.errorHandler(err)
					continue
				}
				if len(batch) == 0 {
					continue
				}

				ctx := modules.NewBatchHandlerContext(batch, ep, i.chainID)
				if err := i.batchHandler(ctx); err != nil {
					i.errorHandler(err)
					continue
				}

				senders := append(getSenders(ctx.Batch), getSenders(ctx.PendingRemoval)...)
				if err := i.mempool.RemoveOps(ep, senders...); err != nil {
					i.errorHandler(err)
					continue
				}
			}

			time.Sleep(5 * time.Second)
		}
	}(i)

	return nil
}
