// Package noop implements basic no-operation modules which are used by default for both Client and Bundler.
package noop

import "github.com/stackup-wallet/stackup-bundler/pkg/modules"

// BatchHandler takes a BatchHandlerCtx and returns nil error.
func BatchHandler(ctx *modules.BatchHandlerCtx) error {
	return nil
}

// UserOpHandler takes a UserOpHandlerCtx and returns nil error.
func UserOpHandler(ctx *modules.UserOpHandlerCtx) error {
	return nil
}
