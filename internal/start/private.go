package start

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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

func runDBGarbageCollection(db *badger.DB) {
	go func(db *badger.DB) {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
		again:
			err := db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		}
	}(db)
}

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

	check := checks.New(rpc, conf.MaxVerificationGas, conf.BundlerCollectorTracer)
	relayer := relay.New(db, eoa, eth, chain, beneficiary, logr)
	paymaster := paymaster.New(db)

	// Init Client
	c := client.New(mem, chain, conf.SupportedEntryPoints, conf.MaxVerificationGas)
	c.SetGetUserOpReceiptFunc(client.GetUserOpReceiptWithEthClient(eth))
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
		check.PaymasterDeposit(),
		relayer.SendUserOperation(),
		paymaster.IncOpsIncluded(),
	)
	if err := b.Run(); err != nil {
		log.Fatal(err)
	}

	// init Debug
	var d *client.Debug
	if conf.DebugMode {
		d = client.NewDebug(eoa, eth, mem, b, chain, conf.SupportedEntryPoints[0], beneficiary)
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
	r.POST(
		"/",
		relayer.FilterByClientID(),
		jsonrpc.Controller(client.NewRpcAdapter(c, d)),
		relayer.MapUserOpHashToClientID(),
	)
	if err := r.Run(fmt.Sprintf(":%d", conf.Port)); err != nil {
		log.Fatal(err)
	}
}
