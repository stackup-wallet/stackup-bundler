// Package bundler provides the mediator for processing outgoing UserOperation batches to the EntryPoint.
package bundler

import (
	"context"
	"math/big"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-logr/logr"
	"github.com/stackup-wallet/stackup-bundler/internal/logger"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/gasprice"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/noop"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/multierr"
)

// Bundler controls the end to end process of creating a batch of UserOperations from the mempool and sending
// it to the EntryPoint.
type Bundler struct {
	mempool              *mempool.Mempool
	chainID              *big.Int
	supportedEntryPoints []common.Address
	batchHandler         modules.BatchHandlerFunc
	logger               logr.Logger
	meter                metric.Meter
	isRunning            bool
	done                 chan bool
	stop                 func()
	maxBatch             int
	gbf                  gasprice.GetBaseFeeFunc
	ggt                  gasprice.GetGasTipFunc
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
		meter:                otel.GetMeterProvider().Meter("bundler"),
		isRunning:            false,
		done:                 make(chan bool),
		stop:                 func() {},
		maxBatch:             0,
		gbf:                  gasprice.NoopGetBaseFeeFunc(),
		ggt:                  gasprice.NoopGetGasTipFunc(),
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

// SetGetGasTipFunc defines the function used to retrieve an estimate for gas tip during each bundler run.
func (i *Bundler) SetGetGasTipFunc(ggt gasprice.GetGasTipFunc) {
	i.ggt = ggt
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

// UserMeter defines an opentelemetry meter object used by the Bundler instance to capture metrics during each
// run.
func (i *Bundler) UserMeter(meter metric.Meter) error {
	i.meter = meter
	_, err := i.meter.Int64ObservableGauge(
		"bundler_mempool_size",
		metric.WithInt64Callback(func(ctx context.Context, io metric.Int64Observer) error {
			size := 0
			for _, ep := range i.supportedEntryPoints {
				batch, err := i.mempool.Dump(ep)
				if err != nil {
					return err
				}
				size += len(batch)
			}
			io.Observe(int64(size))
			return nil
		}),
	)
	return err
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

	// Get all pending userOps from the mempool. This will be in FIFO order. Downstream modules should sort it
	// based on more specific strategies.
	batch, err := i.mempool.Dump(ep)
	if err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}
	if len(batch) == 0 {
		return nil, nil
	}
	batch = adjustBatchSize(i.maxBatch, batch)

	// Get current block basefee
	bf, err := i.gbf()
	if err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}

	// Get suggested gas tip
	var gt *big.Int
	if bf != nil {
		gt, err = i.ggt()
		if err != nil {
			l.Error(err, "bundler run error")
			return nil, err
		}
	}

	// Get suggested gas price (for networks that don't support EIP-1559)
	gp, err := i.ggp()
	if err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}

	// Create context and execute modules.
	ctx := modules.NewBatchHandlerContext(batch, ep, i.chainID, bf, gt, gp)
	if err := i.batchHandler(ctx); err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}

	// Remove userOps that remain in the context from mempool.
	rmOps := append([]*userop.UserOperation{}, ctx.Batch...)
	rmOps = append(rmOps, ctx.PendingRemoval...)
	if err := i.mempool.RemoveOps(ep, rmOps...); err != nil {
		l.Error(err, "bundler run error")
		return nil, err
	}

	// Update logs for the current run.
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

	// Construct an exponential backoff to avoid overwhelming the system if error
	// happens during the ticking.
	bo := backoff.NewExponentialBackOff()
	{
		bo.InitialInterval = 5 * time.Second
		bo.Multiplier = 2
		bo.MaxElapsedTime = 0
		bo.Reset()
	}

	ticker := time.NewTicker(bo.InitialInterval)
	go func(i *Bundler) {
		for {
			select {
			case <-i.done:
				return
			case <-ticker.C:
				if err := i.runOnce(); err != nil {
					// Use exponential backoff.
					ticker.Reset(bo.NextBackOff())
					continue
				}

				// Reset back to normal ticking.
				ticker.Reset(bo.InitialInterval)
				bo.Reset()
			}
		}
	}(i)

	i.isRunning = true
	i.stop = ticker.Stop
	return nil
}

func (i *Bundler) runOnce() error {
	var tickErrs []error
	for _, ep := range i.supportedEntryPoints {
		_, err := i.Process(ep)
		if err != nil {
			// Already logged.
			tickErrs = append(tickErrs, err)
			continue
		}
	}

	return multierr.Combine(tickErrs...)
}

// Stop signals the bundler to stop continuously processing batches from the mempool.
func (i *Bundler) Stop() {
	if !i.isRunning {
		return
	}

	i.isRunning = false
	i.stop()
	i.done <- true
}
