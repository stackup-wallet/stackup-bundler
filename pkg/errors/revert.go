package errors

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func revertError() abi.Error {
	reason, _ := abi.NewType("string", "string", nil)
	return abi.NewError("Error", abi.Arguments{
		{Name: "reason", Type: reason},
	})
}

func DecodeRevert(data []byte) (string, error) {
	abi := revertError()
	revert, err := abi.Unpack(data)
	if err != nil {
		return "", fmt.Errorf("revert: %s", err)
	}

	args, ok := revert.([]any)
	if !ok {
		return "", errors.New("revert: cannot assert type: args is not of type []any")
	}
	if len(args) != 1 {
		return "", fmt.Errorf("revert: invalid args length: expected 1, got %d", len(args))
	}

	return args[0].(string), nil
}
