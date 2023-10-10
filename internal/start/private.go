package start

import (
	"context"
	"fmt"
	"log"
	"net/http"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stackup-wallet/stackup-bundler/internal/config"
	"github.com/stackup-wallet/stackup-bundler/internal/logger"
	"github.com/stackup-wallet/stackup-bundler/internal/o11y"
	"github.com/stackup-wallet/stackup-bundler/pkg/bundler"
	"github.com/stackup-wallet/stackup-bundler/pkg/client"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/jsonrpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/batch"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/checks"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/expire"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/gasprice"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/paymaster"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/relay"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

func PrivateMode() {
	conf := config.GetValues()

	logr := logger.NewZeroLogr().
		WithName("stackup_bundler").
		WithValues("bundler_mode", "private")

	eoa, err := signer.New(conf.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	beneficiary := common.HexToAddress(conf.Beneficiary)

	db, err := badger.Open(badger.DefaultOptions(conf.DataDirectory))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	runDBGarbageCollection(db)

	rpc, err := rpc.Dial(conf.EthClientUrl)
	if err != nil {
		log.Fatal(err)
	}

	eth := ethclient.NewClient(rpc)

	chain, err := eth.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if o11y.IsEnabled(conf.OTELServiceName) {
		o11yOpts := &o11y.Opts{
			ServiceName:     conf.OTELServiceName,
			CollectorHeader: conf.OTELCollectorHeaders,
			CollectorUrl:    conf.OTELCollectorUrl,
			InsecureMode:    conf.OTELInsecureMode,

			ChainID: chain,
			Address: eoa.Address,
		}

		tracerCleanup := o11y.InitTracer(o11yOpts)
		defer tracerCleanup()

		metricsCleanup := o11y.InitMetrics(o11yOpts)
		defer metricsCleanup()
	}

	ov := gas.NewDefaultOverhead()
	if chain.Cmp(config.ArbitrumOneChainID) == 0 || chain.Cmp(config.ArbitrumGoerliChainID) == 0 {
		ov.SetCalcPreVerificationGasFunc(gas.CalcArbitrumPVGWithEthClient(rpc, conf.SupportedEntryPoints[0]))
		ov.SetPreVerificationGasBufferFactor(16)
	}
	if chain.Cmp(config.OptimismChainID) == 0 || chain.Cmp(config.OptimismGoerliChainID) == 0 ||
		chain.Cmp(config.BaseChainID) == 0 || chain.Cmp(config.BaseGoerliChainID) == 0 {
		ov.SetCalcPreVerificationGasFunc(
			gas.CalcOptimismPVGWithEthClient(rpc, chain, conf.SupportedEntryPoints[0]),
		)
		ov.SetPreVerificationGasBufferFactor(1)
	}

	mem, err := mempool.New(db)
	if err != nil {
		log.Fatal(err)
	}

	check := checks.New(
		db,
		rpc,
		ov,
		conf.MaxVerificationGas,
		conf.MaxBatchGasLimit,
		conf.MaxOpsForUnstakedSender,
	)

	exp := expire.New(conf.MaxOpTTL)

	relayer := relay.New(eoa, eth, chain, beneficiary, logr)

	paymaster := paymaster.New(db)

	// Init Client
	c := client.New(mem, ov, chain, conf.SupportedEntryPoints)
	c.SetGetUserOpReceiptFunc(client.GetUserOpReceiptWithEthClient(eth))
	c.SetGetGasEstimateFunc(
		client.GetGasEstimateWithEthClient(rpc, ov, chain, conf.MaxBatchGasLimit, conf.PMGasEstBuffer),
	)
	c.SetGetUserOpByHashFunc(client.GetUserOpByHashWithEthClient(eth))
	c.UseLogger(logr)
	c.UseModules(
		check.ValidateOpValues(),
		paymaster.CheckStatus(),
		check.SimulateOp(),
		paymaster.IncOpsSeen(),
	)

	// Init Bundler
	b := bundler.New(mem, chain, conf.SupportedEntryPoints)
	b.SetGetBaseFeeFunc(gasprice.GetBaseFeeWithEthClient(eth))
	b.SetGetGasTipFunc(gasprice.GetGasTipWithEthClient(eth))
	b.SetGetLegacyGasPriceFunc(gasprice.GetLegacyGasPriceWithEthClient(eth))
	b.SetMaxBatch(conf.MaxOpsForBundle)
	b.UseLogger(logr)
	if err := b.UserMeter(otel.GetMeterProvider().Meter("bundler")); err != nil {
		log.Fatal(err)
	}
	b.UseModules(
		exp.DropExpired(),
		gasprice.SortByGasPrice(),
		gasprice.FilterUnderpriced(),
		batch.SortByNonce(),
		batch.MaintainGasLimit(conf.MaxBatchGasLimit),
		check.CodeHashes(),
		check.PaymasterDeposit(),
		relayer.SendUserOperation(),
		paymaster.IncOpsIncluded(),
		check.Clean(),
	)
	if err := b.Run(); err != nil {
		log.Fatal(err)
	}

	// init Debug
	var d *client.Debug
	if conf.DebugMode {
		d = client.NewDebug(eoa, eth, mem, b, chain, conf.SupportedEntryPoints[0], beneficiary)
		b.SetMaxBatch(1)
		relayer.SetWaitTimeout(0)
	}

	// Init HTTP server
	gin.SetMode(conf.GinMode)
	r := gin.New()
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}
	if o11y.IsEnabled(conf.OTELServiceName) {
		r.Use(otelgin.Middleware(conf.OTELServiceName))
	}
	r.Use(
		cors.Default(),
		logger.WithLogr(logr),
		gin.Recovery(),
	)
	r.GET("/ping", func(g *gin.Context) {
		g.Status(http.StatusOK)
	})
	handlers := []gin.HandlerFunc{
		jsonrpc.Controller(client.NewRpcAdapter(c, d)),
		jsonrpc.WithOTELTracerAttributes(),
	}
	r.POST("/", handlers...)
	r.POST("/rpc", handlers...)

	if err := r.Run(fmt.Sprintf(":%d", conf.Port)); err != nil {
		log.Fatal(err)
	}
}
