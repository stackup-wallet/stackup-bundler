package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/stackup-wallet/stackup-bundler/internal/config"
	"github.com/stackup-wallet/stackup-bundler/internal/jsonrpc"
	"github.com/stackup-wallet/stackup-bundler/internal/redispool"
	"github.com/stackup-wallet/stackup-bundler/pkg/client"
)

func main() {
	conf := config.GetValues()

	eth, err := ethclient.Dial(conf.RpcUrl)
	if err != nil {
		log.Fatal(err)
	}
	mem, err := redispool.NewClientInterface(conf.RedisUrl)
	if err != nil {
		log.Fatal(err)
	}
	client := client.New(eth, mem, conf.SupportedEntryPoints)

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
