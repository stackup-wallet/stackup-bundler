package errors

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func panicError() abi.Error {
	code, _ := abi.NewType("uint256", "uint256", nil)
	return abi.NewError("Panic", abi.Arguments{
		{Name: "code", Type: code},
	})
}

func DecodePanic(data []byte) (string, error) {
	abi := panicError()
	panic, err := abi.Unpack(data)
	if err != nil {
		return "", fmt.Errorf("panic: %s", err)
	}

	args, ok := panic.([]any)
	if !ok {
		return "", errors.New("panic: cannot assert type: args is not of type []any")
	}
	if len(args) != 1 {
		return "", fmt.Errorf("panic: invalid args length: expected 1, got %d", len(args))
	}

	return hexutil.EncodeBig(args[0].(*big.Int)), nil
}
