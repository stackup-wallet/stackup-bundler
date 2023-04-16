package reverts

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

type ExecutionResultRevert struct {
	PreOpGas      *big.Int
	Paid          *big.Int
	ValidAfter    *big.Int
	ValidUntil    *big.Int
	TargetSuccess bool
	TargetResult  []byte
}

func executionResult() abi.Error {
	uint256, _ := abi.NewType("uint256", "", nil)
	uint48, _ := abi.NewType("uint48", "", nil)
	boolean, _ := abi.NewType("bool", "", nil)
	bytes, _ := abi.NewType("bytes", "", nil)
	return abi.NewError("ExecutionResult", abi.Arguments{
		{Name: "preOpGas", Type: uint256},
		{Name: "paid", Type: uint256},
		{Name: "validAfter", Type: uint48},
		{Name: "validUntil", Type: uint48},
		{Name: "targetSuccess", Type: boolean},
		{Name: "targetResult", Type: bytes},
	})
}

func NewExecutionResult(err error) (*ExecutionResultRevert, error) {
	rpcErr, ok := err.(rpc.DataError)
	if !ok {
		return nil, errors.New("executionResult: cannot assert type: error is not of type rpc.DataError")
	}

	data, ok := rpcErr.ErrorData().(string)
	if !ok {
		return nil, errors.New("executionResult: cannot assert type: data is not of type string")
	}

	sim := executionResult()
	revert, err := sim.Unpack(common.Hex2Bytes(data[2:]))
	if err != nil {
		return nil, fmt.Errorf("executionResult: %s", err)
	}

	args, ok := revert.([]any)
	if !ok {
		return nil, errors.New("executionResult: cannot assert type: args is not of type []any")
	}
	if len(args) != 6 {
		return nil, fmt.Errorf("executionResult: invalid args length: expected 6, got %d", len(args))
	}

	return &ExecutionResultRevert{
		PreOpGas:      args[0].(*big.Int),
		Paid:          args[1].(*big.Int),
		ValidAfter:    args[2].(*big.Int),
		ValidUntil:    args[3].(*big.Int),
		TargetSuccess: args[4].(bool),
		TargetResult:  args[5].([]byte),
	}, nil
}
