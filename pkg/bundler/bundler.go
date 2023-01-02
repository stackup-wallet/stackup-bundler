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
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Bundler controls the end to end process of creating a batch of UserOperations from the mempool and sending
// it to the EntryPoint.
type Bundler struct {
	mempool              *mempool.Mempool
	chainID              *big.Int
	supportedEntryPoints []common.Address
	batchHandler         modules.BatchHandlerFunc
	logger               logr.Logger
	isRunning            bool
	stop                 chan bool
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
		isRunning:            false,
		stop:                 make(chan bool),
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
	if i.isRunning {
		return nil
	}

	go func(i *Bundler) {
		logger := i.logger.WithName("run")

		for {
			select {
			case <-i.stop:
				return
			default:
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

					rmOps := append([]*userop.UserOperation{}, ctx.Batch...)
					rmOps = append(rmOps, ctx.PendingRemoval...)
					if err := i.mempool.RemoveOps(ep, rmOps...); err != nil {
						l.Error(err, "bundler run error")
						continue
					}

					bat := []string{}
					for _, op := range ctx.Batch {
						bat = append(bat, op.GetUserOpHash(ep, i.chainID).String())
					}
					l = l.WithValues("batch_userop_hashes", bat)

					drp := []string{}
					for _, op := range ctx.PendingRemoval {
						drp = append(drp, op.GetUserOpHash(ep, i.chainID).String())
					}
					l = l.WithValues("dropped_userop_hashes", drp)

					for k, v := range ctx.Data {
						l = l.WithValues(k, v)
					}
					l = l.WithValues("duration", time.Since(start))
					l.Info("bundler run ok")
				}
			}
		}
	}(i)

	i.isRunning = true
	return nil
}

// Stop signals the bundler to stop continuously processing batches from the mempool.
func (i *Bundler) Stop() {
	if !i.isRunning {
		return
	}

	i.isRunning = false
	i.stop <- true
}
