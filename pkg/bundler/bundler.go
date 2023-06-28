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
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/gasprice"
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
	watch                chan bool
	onStop               func()
	maxBatch             int
	gbf                  gasprice.GetBaseFeeFunc
	ggp                  gasprice.GetLegacyGasPriceFunc
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
		watch:                make(chan bool),
		onStop:               func() {},
		maxBatch:             0,
		gbf:                  gasprice.NoopGetBaseFeeFunc(),
		ggp:                  gasprice.NoopGetLegacyGasPriceFunc(),
	}
}

// SetMaxBatch defines the max number of UserOperations per bundle. The default value is 0 (i.e. unlimited).
func (i *Bundler) SetMaxBatch(max int) {
	i.maxBatch = max
}

// SetGetBaseFeeFunc defines the function used to retrieve an estimate for basefee during each bundler run.
func (i *Bundler) SetGetBaseFeeFunc(gbf gasprice.GetBaseFeeFunc) {
	i.gbf = gbf
}

// SetGetLegacyGasPriceFunc defines the function used to retrieve an estimate for gas price during each
// bundler run.
func (i *Bundler) SetGetLegacyGasPriceFunc(ggp gasprice.GetLegacyGasPriceFunc) {
	i.ggp = ggp
}

// UseLogger defines the logger object used by the Bundler instance based on the go-logr/logr interface.
func (i *Bundler) UseLogger(logger logr.Logger) {
	i.logger = logger.WithName("bundler")
}

// UseModules defines the BatchHandlers to process batches after it has gone through the standard checks.
func (i *Bundler) UseModules(handlers ...modules.BatchHandlerFunc) {
	i.batchHandler = modules.ComposeBatchHandlerFunc(handlers...)
}

// Process will create a batch from the mempool and send it through to the EntryPoint.
func (i *Bundler) Process(ep common.Address) (*modules.BatchHandlerCtx, error) {
	// Init logger
	start := time.Now()
	l := i.logger.
		WithName("run").
		WithValues("entrypoint", ep.String()).
		WithValues("chain_id", i.chainID.String())

	// Get current block basefee
	bf, err := i.gbf()
	if err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}

	// Get suggested gas price (for networks that don't support EIP-1559)
	gp, err := i.ggp()
	if err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}

	// Get all pending userOps from the mempool. This will be in FIFO order. Downstream modules should sort it
	// accordingly.
	batch, err := i.mempool.BundleOps(ep)
	if err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}
	if len(batch) == 0 {
		return nil, nil
	}
	batch = adjustBatchSize(i.maxBatch, batch)

	// Create context and execute modules.
	ctx := modules.NewBatchHandlerContext(batch, ep, i.chainID, bf, gp)
	if err := i.batchHandler(ctx); err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}

	// Remove included and dropped userOps from the mempool.
	rmOps := append([]*userop.UserOperation{}, ctx.Batch...)
	rmOps = append(rmOps, ctx.PendingRemoval...)
	if err := i.mempool.RemoveOps(ep, rmOps...); err != nil {
		l.Error(err, "bundler run error")
		return nil, err
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
	return ctx, nil
}

// Run starts a goroutine that will continuously process batches from the mempool.
func (i *Bundler) Run() error {
	if i.isRunning {
		return nil
	}

	go func(i *Bundler) {
		for {
			select {
			case <-i.stop:
				return
			case <-i.watch:
				for _, ep := range i.supportedEntryPoints {
					_, err := i.Process(ep)
					if err != nil {
						// Already logged.
						continue
					}
				}
			}
		}
	}(i)

	i.isRunning = true
	i.onStop = i.mempool.OnAdd(i.watch)
	return nil
}

// Stop signals the bundler to stop continuously processing batches from the mempool.
func (i *Bundler) Stop() {
	if !i.isRunning {
		return
	}

	i.isRunning = false
	i.stop <- true
	i.onStop()
}
