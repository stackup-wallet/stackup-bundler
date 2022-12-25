// Package modules provides standard interfaces for extending the Client and Bundler packages with
// middleware.
package modules

// BatchHandlerFunc is an interface to support modular processing of UserOperation batches by the Bundler.
type BatchHandlerFunc func(ctx *BatchHandlerCtx) error

// OpHandlerFunc is an interface to support modular processing of single UserOperations by the Client.
type UserOpHandlerFunc func(ctx *UserOpHandlerCtx) error
