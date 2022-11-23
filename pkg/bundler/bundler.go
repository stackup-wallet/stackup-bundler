// Package bundler provides the mediator for processing outgoing UserOperation batches to the EntryPoint.
package bundler

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-logr/logr"
	"github.com/stackup-wallet/stackup-bundler/internal/logger"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/noop"
)

// Bundler controls the end to end process of creating a batch of UserOperations from the mempool and sending
// it to the EntryPoint.
type Bundler struct {
	mempool              *mempool.Mempool
	chainID              *big.Int
	supportedEntryPoints []common.Address
	batchHandler         modules.BatchHandlerFunc
	logger               logr.Logger
}

// New initializes a new EIP-4337 bundler which can be extended with modules for validating batches and
// excluding UserOperations that should not be sent to the EntryPoint and/or dropped from the mempool.
func New(mempool *mempool.Mempool, chainID *big.Int, supportedEntryPoints []common.Address) *Bundler {
	return &Bundler{
		mempool:              mempool,
		chainID:              chainID,
		supportedEntryPoints: supportedEntryPoints,
		batchHandler:         noop.BatchHandler,
		logger:               logger.NewZeroLogr().WithName("bundler"),
	}
}

// UseLogger defines the logger object used by the Bundler instance based on the go-logr/logr interface.
func (i *Bundler) UseLogger(logger logr.Logger) {
	i.logger = logger.WithName("bundler")
}

// UseModules defines the BatchHandlers to process batches after it has gone through the standard checks.
func (i *Bundler) UseModules(handlers ...modules.BatchHandlerFunc) {
	i.batchHandler = modules.ComposeBatchHandlerFunc(handlers...)
}

// Run starts a goroutine that will continuously process batches from the mempool.
func (i *Bundler) Run() error {
	go func(i *Bundler) {
		logger := i.logger.WithName("run")

		for {
			for _, ep := range i.supportedEntryPoints {
				start := time.Now()
				l := logger.
					WithValues("entrypoint", ep.String()).
					WithValues("chain_id", i.chainID.String())

				batch, err := i.mempool.BundleOps(ep)
				if err != nil {
					l.Error(err, "bundler run error")
					continue
				}
				if len(batch) == 0 {
					continue
				}

				ctx := modules.NewBatchHandlerContext(batch, ep, i.chainID)
				if err := i.batchHandler(ctx); err != nil {
					l.Error(err, "bundler run error")
					continue
				}

				senders := append(getSenders(ctx.Batch), getSenders(ctx.PendingRemoval)...)
				if err := i.mempool.RemoveOps(ep, senders...); err != nil {
					l.Error(err, "bundler run error")
					continue
				}

				bat := []string{}
				for _, op := range ctx.Batch {
					bat = append(bat, op.GetRequestID(ep, i.chainID).String())
				}
				l = l.WithValues("batch_request_ids", bat)

				drp := []string{}
				for _, op := range ctx.PendingRemoval {
					drp = append(drp, op.GetRequestID(ep, i.chainID).String())
				}
				l = l.WithValues("dropped_request_ids", drp)

				for k, v := range ctx.Data {
					l = l.WithValues(k, v)
				}
				l = l.WithValues("duration", time.Since(start))
				l.Info("bundler run ok")
			}
		}
	}(i)

	return nil
}
