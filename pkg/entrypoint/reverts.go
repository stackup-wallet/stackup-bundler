package entrypoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

type PaymasterInfo struct {
	PaymasterStake        *big.Int `json:"paymasterStake"`
	PaymasterUnstakeDelay *big.Int `json:"paymasterUnstakeDelay"`
}

type SimulationResultRevert struct {
	PreOpGas      *big.Int
	Prefund       *big.Int
	Deadline      *big.Int
	PaymasterInfo *PaymasterInfo
}

func simulationResult() abi.Error {
	uint256, _ := abi.NewType("uint256", "uint256", nil)
	paymasterInfo, _ := abi.NewType("tuple", "PaymasterInfo", []abi.ArgumentMarshaling{
		{Name: "paymasterStake", Type: "uint256"},
		{Name: "paymasterUnstakeDelay", Type: "uint256"},
	})
	return abi.NewError("SimulationResult", abi.Arguments{
		{Name: "preOpGas", Type: uint256},
		{Name: "prefund", Type: uint256},
		{Name: "deadline", Type: uint256},
		{Name: "paymasterInfo", Type: paymasterInfo},
	})
}

func newSimulationResultRevert(err error) (*SimulationResultRevert, error) {
	rpcErr, ok := err.(rpc.DataError)
	if !ok {
		return nil, errors.New("simulationResult: cannot assert type: error is not of type rpc.DataError")
	}

	data, ok := rpcErr.ErrorData().(string)
	if !ok {
		return nil, errors.New("simulationResult: cannot assert type: data is not of type string")
	}

	sim := simulationResult()
	revert, err := sim.Unpack(common.Hex2Bytes(data[2:]))
	if err != nil {
		return nil, fmt.Errorf("simulationResult: %s", err)
	}

	args, ok := revert.([]any)
	if !ok {
		return nil, errors.New("simulationResult: cannot assert type: args is not of type []any")
	}
	if len(args) != 4 {
		return nil, fmt.Errorf("simulationResult: invalid args length: expected 4, got %d", len(args))
	}

	preOpGas, ok := args[0].(*big.Int)
	if !ok {
		return nil, errors.New("simulationResult: cannot assert type: preOpGas is not of type *big.Int")
	}

	prefund, ok := args[1].(*big.Int)
	if !ok {
		return nil, errors.New("simulationResult: cannot assert type: prefund is not of type *big.Int")
	}

	deadline, ok := args[2].(*big.Int)
	if !ok {
		return nil, errors.New("simulationResult: cannot assert type: deadline is not of type *big.Int")
	}

	pmi, err := json.Marshal(args[3])
	if err != nil {
		return nil, fmt.Errorf("simulationResult: %s", err)
	}

	paymasterInfo := &PaymasterInfo{}
	if err := json.Unmarshal(pmi, paymasterInfo); err != nil {
		return nil, fmt.Errorf("simulationResult: %s", err)
	}

	return &SimulationResultRevert{
		PreOpGas:      preOpGas,
		Prefund:       prefund,
		Deadline:      deadline,
		PaymasterInfo: paymasterInfo,
	}, nil
}

type FailedOpRevert struct {
	OpIndex   int
	Paymaster common.Address
	Reason    string
}

func failedOp() abi.Error {
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
		return nil, errors.New("failedOp: cannot assert type: error is not of type rpc.DataError")
	}

	data, ok := rpcErr.ErrorData().(string)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: data is not of type string")
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
	if len(args) != 3 {
		return nil, fmt.Errorf("failedOp: invalid args length: expected 3, got %d", len(args))
	}

	opIndex, ok := args[0].(*big.Int)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: opIndex is not of type *big.Int")
	}

	paymaster, ok := args[1].(common.Address)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: paymaster is not of type common.Address")
	}

	reason, ok := args[2].(string)
	if !ok {
		return nil, errors.New("failedOp: cannot assert type: reason is not of type string")
	}

	return &FailedOpRevert{
		OpIndex:   int(opIndex.Int64()),
		Paymaster: paymaster,
		Reason:    reason,
	}, nil
}
