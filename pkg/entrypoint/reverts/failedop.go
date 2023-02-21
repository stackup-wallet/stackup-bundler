package reverts

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

type FailedOpRevert struct {
	OpIndex int
	Reason  string
}

func failedOp() abi.Error {
	opIndex, _ := abi.NewType("uint256", "uint256", nil)
	reason, _ := abi.NewType("string", "string", nil)
	return abi.NewError("FailedOp", abi.Arguments{
		{Name: "opIndex", Type: opIndex},
		{Name: "reason", Type: reason},
	})
}

func NewFailedOp(err error) (*FailedOpRevert, error) {
	rpcErr, ok := err.(rpc.DataError)
	if !ok {
		return nil, fmt.Errorf(
			"failedOp: cannot assert type: error is not of type rpc.DataError, err: %s",
			err,
		)
	}

	data, ok := rpcErr.ErrorData().(string)
	if !ok {
		return nil, fmt.Errorf(
			"failedOp: cannot assert type: data is not of type string, err: %s, data: %s",
			rpcErr.Error(),
			rpcErr.ErrorData(),
		)
	}

	failedOp := failedOp()
	revert, err := failedOp.Unpack(common.Hex2Bytes(data[2:]))
	if err != nil {
		return nil, fmt.Errorf("failedOp: %s", err)
	}

	args, ok := revert.([]any)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: args is not of type []any")
	}
	if len(args) != 2 {
		return nil, fmt.Errorf("failedOp: invalid args length: expected 2, got %d", len(args))
	}

	opIndex, ok := args[0].(*big.Int)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: opIndex is not of type *big.Int")
	}

	reason, ok := args[1].(string)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: reason is not of type string")
	}

	return &FailedOpRevert{
		OpIndex: int(opIndex.Int64()),
		Reason:  reason,
	}, nil
}
