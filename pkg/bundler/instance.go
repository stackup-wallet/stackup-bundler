package bundler

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
)

// Instance is a representation of an ERC-4337 bundler.
type Instance struct {
	mempool              *mempool.Interface
	chainID              *big.Int
	supportedEntryPoints []common.Address
	batchHandler         modules.BatchHandlerFunc
	errorHandler         modules.ErrorHandlerFunc
}

// UseModules defines the BatchHandlers to process batches after it has gone through the standard checks.
func (i *Instance) UseModules(handlers ...modules.BatchHandlerFunc) {
	i.batchHandler = modules.ComposeBatchHandlerFunc(handlers...)
}

// SetErrorHandlerFunc defines a method for handling errors at any point of the process.
func (i *Instance) SetErrorHandlerFunc(handler modules.ErrorHandlerFunc) {
	i.errorHandler = handler
}

// Run starts a goroutine that will continuously process batches from the mempool.
func (i *Instance) Run() error {
	go func(i *Instance) {
		for {
			for _, ep := range i.supportedEntryPoints {
				batch, err := i.mempool.BundleOps(ep)
				if err != nil {
					i.errorHandler(err)
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
