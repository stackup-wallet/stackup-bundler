package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/stackup-wallet/stackup-bundler/internal/config"
	"github.com/stackup-wallet/stackup-bundler/pkg/client"
	"github.com/stackup-wallet/stackup-bundler/pkg/jsonrpc"
)

func main() {
	v := config.GetValues()
	c, err := ethclient.Dial(v.RpcUrl)
	if err != nil {
		log.Fatal(err)
	}
	i := client.New(c, v.SupportedEntryPoints)

	gin.SetMode(v.GinMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	// JSON-RPC handler
	r.POST("/", func(c *gin.Context) { jsonrpc.HandleRequest(c, &i) })

	r.Run(fmt.Sprintf(":%d", v.Port))
}
