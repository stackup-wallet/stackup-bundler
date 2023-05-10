package errors

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

type hexDataErrorWrapper struct {
	output string
}

// ParseHexToRpcDataError is a utility function that converts a hex string into an interface that's compatible
// with rpc.DataError. This is useful for parsing output from a debug_traceCall that resulted in a transaction
// revert.
func ParseHexToRpcDataError(hex string) (rpc.DataError, error) {
	if _, err := common.ParseHexOrString(hex); err != nil {
		return nil, err
	}

	res := &hexDataErrorWrapper{output: hex}
	return res, nil
}

func (t *hexDataErrorWrapper) Error() string {
	return t.output
}
func (t *hexDataErrorWrapper) ErrorData() any {
	return t.output
}
