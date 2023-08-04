package modules

import "fmt"

// ComposeBatchHandlerFunc combines many BatchHandlers into one.
func ComposeBatchHandlerFunc(fns ...BatchHandlerFunc) BatchHandlerFunc {
	return func(ctx *BatchHandlerCtx) error {
		for i, fn := range fns {
			err := fn(ctx)
			if err != nil {
				return fmt.Errorf("error at batch handler %d: %w", i, err)
			}
		}

		return nil
	}
}

// ComposeUserOpHandlerFunc combines many UserOpHandlers into one.
func ComposeUserOpHandlerFunc(fns ...UserOpHandlerFunc) UserOpHandlerFunc {
	return func(ctx *UserOpHandlerCtx) error {
		for _, fn := range fns {
			err := fn(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
