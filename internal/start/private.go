package start

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/stackup-wallet/stackup-bundler/internal/config"
	"github.com/stackup-wallet/stackup-bundler/pkg/bundler"
	"github.com/stackup-wallet/stackup-bundler/pkg/client"
	"github.com/stackup-wallet/stackup-bundler/pkg/jsonrpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/checks"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/println"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/relay"
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

	db, err := badger.Open(badger.DefaultOptions(conf.DataDirectory))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	runDBGarbageCollection(db)

	eth, err := ethclient.Dial(conf.EthClientUrl)
	if err != nil {
		log.Fatal(err)
	}

	chain, err := eth.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	mem, err := mempool.NewBadgerDBWrapper(db)
	if err != nil {
		log.Fatal(err)
	}

	relayer := relay.New(db)
	relayer.SetErrorHandlerFunc(println.ErrorHandler)

	// Start bundler
	b := bundler.New(mem, chain, conf.SupportedEntryPoints)
	b.UseModules(
		println.BatchHandler,
		relayer.SendUserOperation(),
	)
	b.SetErrorHandlerFunc(println.ErrorHandler)
	b.Run()

	// Start client
	c := client.New(mem, chain, conf.SupportedEntryPoints)
	c.UseModules(
		checks.StandaloneClient(eth, mem),
		println.UserOpHandler,
	)

	gin.SetMode(conf.GinMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.GET("/ping", func(g *gin.Context) {
		g.Status(http.StatusOK)
	})
	r.POST(
		"/",
		relayer.FilterByClient(),
		jsonrpc.Controller(client.NewRpcAdapter(c)),
		relayer.MapRequestIDToClientID(chain),
	)
	r.Run(fmt.Sprintf(":%d", conf.Port))
}
