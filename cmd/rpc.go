package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jsonrpc "github.com/stackup-wallet/standalone-bundler/pkg/handler"
)

type TestRPC struct {
	counter int
}

// eth_add test method
func (t *TestRPC) EthAdd(arg int) int {
	t.counter += arg
	return t.counter
}

// eth_sub test method
func (t *TestRPC) EthSub(arg int) int {
	t.counter -= arg
	return t.counter
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// TODO: Implement ERC-4337 Client interface
	rpc := TestRPC{}

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	// JSON-RPC handler
	r.POST("/", func(c *gin.Context) { jsonrpc.HandleRequest(c, &rpc) })

	return r
}

func main() {
	r := setupRouter()

	// Listen and Server in 0.0.0.0:4337
	r.Run(":4337")
}
