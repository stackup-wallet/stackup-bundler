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

type StakeInfo struct {
	Stake           *big.Int `json:"stake"`
	UnstakeDelaySec *big.Int `json:"unstakeDelaySec"`
}

type SimulationResultRevert struct {
	PreOpGas      *big.Int
	Prefund       *big.Int
	Deadline      *big.Int
	SenderInfo    *StakeInfo
	FactoryInfo   *StakeInfo
	PaymasterInfo *StakeInfo
}

var (
	stakeInfoType = []abi.ArgumentMarshaling{
		{Name: "stake", Type: "uint256"},
		{Name: "unstakeDelaySec", Type: "uint256"},
	}
)

func simulationResult() abi.Error {
	uint256, _ := abi.NewType("uint256", "uint256", nil)
	senderInfo, _ := abi.NewType("tuple", "SenderInfo", stakeInfoType)
	factoryInfo, _ := abi.NewType("tuple", "FactoryInfo", stakeInfoType)
	paymasterInfo, _ := abi.NewType("tuple", "PaymasterInfo", stakeInfoType)

	return abi.NewError("SimulationResult", abi.Arguments{
		{Name: "preOpGas", Type: uint256},
		{Name: "prefund", Type: uint256},
		{Name: "deadline", Type: uint256},
		{Name: "senderInfo", Type: senderInfo},
		{Name: "factoryInfo", Type: factoryInfo},
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
	if len(args) != 6 {
		return nil, fmt.Errorf("simulationResult: invalid args length: expected 6, got %d", len(args))
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

	senderInfo := &StakeInfo{}
	si, err := json.Marshal(args[3])
	if err != nil {
		return nil, fmt.Errorf("simulationResult: %s", err)
	}
	if err := json.Unmarshal(si, senderInfo); err != nil {
		return nil, fmt.Errorf("simulationResult: %s", err)
	}

	factoryInfo := &StakeInfo{}
	fi, err := json.Marshal(args[4])
	if err != nil {
		return nil, fmt.Errorf("simulationResult: %s", err)
	}
	if err := json.Unmarshal(fi, factoryInfo); err != nil {
		return nil, fmt.Errorf("simulationResult: %s", err)
	}

	paymasterInfo := &StakeInfo{}
	pmi, err := json.Marshal(args[5])
	if err != nil {
		return nil, fmt.Errorf("simulationResult: %s", err)
	}
	if err := json.Unmarshal(pmi, paymasterInfo); err != nil {
		return nil, fmt.Errorf("simulationResult: %s", err)
	}

	return &SimulationResultRevert{
		PreOpGas:      preOpGas,
		Prefund:       prefund,
		Deadline:      deadline,
		SenderInfo:    senderInfo,
		FactoryInfo:   factoryInfo,
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
		return nil, fmt.Errorf("failedOp: cannot assert type: error is not of type rpc.DataError, err: %s", err)
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
