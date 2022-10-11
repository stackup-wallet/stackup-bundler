package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stackup-wallet/stackup-bundler/internal/config"
	"github.com/stackup-wallet/stackup-bundler/pkg/client"
	"github.com/stackup-wallet/stackup-bundler/pkg/jsonrpc"
)

func main() {
	v := config.GetValues()
	i := client.New(v.SupportedEntryPoints)

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
