package entrypoint

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

type FailedOpRevert struct {
	OpIndex   int
	Paymaster common.Address
	Reason    string
}

func failedOpError() abi.Error {
	opIndex, _ := abi.NewType("uint256", "uint256", nil)
	paymaster, _ := abi.NewType("address", "address", nil)
	reason, _ := abi.NewType("string", "string", nil)
	return abi.NewError("FailedOp", abi.Arguments{
		{Name: "opIndex", Type: opIndex},
		{Name: "paymaster", Type: paymaster},
		{Name: "reason", Type: reason},
	})
}

func newFailedOpRevert(err error) (*FailedOpRevert, error) {
	rpcErr, ok := err.(rpc.DataError)
	if !ok {
		return nil, errors.New("cannot assert type: error is not of type rpc.DataError")
	}

	data, ok := rpcErr.ErrorData().(string)
	if !ok {
		return nil, errors.New("cannot assert type: data is not of type string")
	}

	failedOp := failedOpError()
	revert, err := failedOp.Unpack(common.Hex2Bytes(data[2:]))
	if err != nil {
		return nil, err
	}

	args, ok := revert.([]any)
	if !ok {
		return nil, errors.New("cannot assert type: args is not of type []any")
	}
	if len(args) != 3 {
		return nil, fmt.Errorf("invalid args length: expected 3, got %d", len(args))
	}

	opIndex, ok := args[0].(*big.Int)
	if !ok {
		return nil, errors.New("cannot assert type: opIndex is not of type *big.Int")
	}

	paymaster, ok := args[1].(common.Address)
	if !ok {
		return nil, errors.New("cannot assert type: paymaster is not of type common.Address")
	}

	reason, ok := args[2].(string)
	if !ok {
		return nil, errors.New("cannot assert type: reason is not of type string")
	}

	return &FailedOpRevert{OpIndex: int(opIndex.Int64()), Paymaster: paymaster, Reason: reason}, nil
}
