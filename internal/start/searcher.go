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
	"github.com/metachris/flashbotsrpc"
	"github.com/stackup-wallet/stackup-bundler/internal/config"
	"github.com/stackup-wallet/stackup-bundler/internal/logger"
	"github.com/stackup-wallet/stackup-bundler/pkg/bundler"
	"github.com/stackup-wallet/stackup-bundler/pkg/client"
	"github.com/stackup-wallet/stackup-bundler/pkg/jsonrpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/builder"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/checks"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/paymaster"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
)

func SearcherMode() {
	conf := config.GetValues()

	logr := logger.NewZeroLogr().
		WithName("stackup_bundler").
		WithValues("bundler_mode", "searcher")

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

	fb := flashbotsrpc.NewFlashbotsRPC(conf.EthBuilderUrl)

	chain, err := eth.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if !builder.CompatibleChainIDs.Contains(chain.Uint64()) {
		log.Fatalf(
			"error: network with chainID %d is not compatible with the Block Builder API.",
			chain.Uint64(),
		)
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
	// TODO: Create separate go-routine for tracking transactions sent to the block builder.
	builder := builder.New(eoa, eth, fb, beneficiary, conf.BlocksInTheFuture)
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
		// TODO: add p2p propagation module
		paymaster.IncOpsSeen(),
	)

	// Init Bundler
	b := bundler.New(mem, chain, conf.SupportedEntryPoints)
	b.UseLogger(logr)
	b.UseModules(
		check.CodeHashes(),
		check.PaymasterDeposit(),
		builder.SendUserOperation(),
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
		jsonrpc.Controller(client.NewRpcAdapter(c, d)),
	}
	r.POST("/", handlers...)
	r.POST("/rpc", handlers...)
	if err := r.Run(fmt.Sprintf(":%d", conf.Port)); err != nil {
		log.Fatal(err)
	}
}
