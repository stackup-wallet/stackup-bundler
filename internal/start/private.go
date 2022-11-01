package start

import (
	"fmt"
	"log"
	"net/http"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/stackup-wallet/stackup-bundler/internal/config"
	"github.com/stackup-wallet/stackup-bundler/internal/jsonrpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/bundler"
	"github.com/stackup-wallet/stackup-bundler/pkg/client"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
)

func startBundler(conf *config.Values, eth *ethclient.Client, mem *mempool.Interface) {
	// Spawn a bundler thread for each supported entryPoint
	for _, ep := range conf.SupportedEntryPoints {
		b, err := bundler.New(eth, mem, ep)
		if err != nil {
			log.Fatal(err)
		}

		go func(b *bundler.Instance) {
			for {
				res, err := b.ProcessBatch()
				if err != nil {
					fmt.Println(fmt.Errorf("bundler error: %w", err))
				} else if !res {
					fmt.Println("Nothing to process")
				}

				time.Sleep(5 * time.Second)
			}
		}(b)
	}

}

func startClient(conf *config.Values, eth *ethclient.Client, mem *mempool.Interface) {
	client, err := client.New(eth, mem, conf.SupportedEntryPoints)
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode(conf.GinMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	// JSON-RPC handler
	r.POST("/", func(c *gin.Context) { jsonrpc.HandleRequest(c, client) })

	r.Run(fmt.Sprintf(":%d", conf.Port))
}

func PrivateMode() {
	conf := config.GetValues()

	db, err := badger.Open(badger.DefaultOptions(conf.DataDirectory))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	eth, err := ethclient.Dial(conf.RpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	mem, err := mempool.NewBadgerDBWrapper(db)
	if err != nil {
		log.Fatal(err)
	}

	startBundler(conf, eth, mem)
	startClient(conf, eth, mem)
}
