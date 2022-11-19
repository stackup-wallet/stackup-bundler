package println

import (
	"fmt"
	"log"

	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
)

// BatchHandler prints each op in the batch to a new line.
func BatchHandler(ctx *modules.BatchHandlerCtx) error {
	fmt.Println("log batch:")
	for _, op := range ctx.Batch {
		b, _ := op.MarshalJSON()
		fmt.Println(string(b))
	}

	return nil
}

// UserOpHandler prints the op to a new line.
func UserOpHandler(ctx *modules.UserOpHandlerCtx) error {
	op, _ := ctx.UserOp.MarshalJSON()
	fmt.Printf("log userOp: %s\n", string(op))
	return nil
}

// ErrorHandler is a simple wrapper around log.Fatal()
func ErrorHandler(err error) {
	log.Fatal(err)
}
