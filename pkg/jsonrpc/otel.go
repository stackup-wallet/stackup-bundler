package jsonrpc

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// WithOTELTracerAttributes adds custom opentelemetry attributes relating to the JSON-RPC method call for the
// current span.
func WithOTELTracerAttributes() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req, ok := ctx.Get("json-rpc-request")
		if ok {
			json := req.(map[string]any)
			span := trace.SpanFromContext(ctx.Request.Context())
			span.SetAttributes(attribute.String("jsonrpc_method", json["method"].(string)))
		}
	}
}
