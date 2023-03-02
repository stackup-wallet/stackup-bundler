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
	"github.com/stackup-wallet/stackup-bundler/pkg/bundler"
	"github.com/stackup-wallet/stackup-bundler/pkg/client"
	"github.com/stackup-wallet/stackup-bundler/pkg/jsonrpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/checks"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/paymaster"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/relay"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
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

	mem, err := mempool.New(db)
	if err != nil {
		log.Fatal(err)
	}

	check := checks.New(
		db,
		rpc,
		conf.MaxVerificationGas,
		conf.MaxOpsForUnstakedSender,
		conf.BundlerCollectorTracer,
	)

	relayer := relay.New(db, eoa, eth, chain, beneficiary, logr)
	if conf.RelayerBannedThreshold > 0 {
		relayer.SetBannedThreshold(conf.RelayerBannedThreshold)
	}
	if conf.RelayerBannedTimeWindow > 0 {
		relayer.SetBannedTimeWindow(conf.RelayerBannedTimeWindow)
	}

	paymaster := paymaster.New(db)

	// Init Client
	c := client.New(mem, chain, conf.SupportedEntryPoints)
	c.SetGetUserOpReceiptFunc(client.GetUserOpReceiptWithEthClient(eth))
	c.SetGetSimulateValidationFunc(client.GetSimulateValidationWithRpcClient(rpc))
	c.SetGetCallGasEstimateFunc(client.GetCallGasEstimateWithEthClient(eth))
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
	b.UseLogger(logr)
	b.UseModules(
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
		relayer.SetBannedThreshold(relay.NoBanThreshold)
	}

	// Init HTTP server
	gin.SetMode(conf.GinMode)
	r := gin.New()
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
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
		relayer.FilterByClientID(),
		jsonrpc.Controller(client.NewRpcAdapter(c, d)),
		relayer.MapUserOpHashToClientID(),
	}
	r.POST("/", handlers...)
	r.POST("/rpc", handlers...)

	if err := r.Run(fmt.Sprintf(":%d", conf.Port)); err != nil {
		log.Fatal(err)
	}
}
