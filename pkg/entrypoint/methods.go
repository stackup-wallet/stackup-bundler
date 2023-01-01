package entrypoint

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	userOp, _ = abi.NewType("tuple", "userOp", []abi.ArgumentMarshaling{
		{Name: "Sender", Type: "address"},
		{Name: "Nonce", Type: "uint256"},
		{Name: "InitCode", Type: "bytes"},
		{Name: "CallData", Type: "bytes"},
		{Name: "CallGasLimit", Type: "uint256"},
		{Name: "VerificationGasLimit", Type: "uint256"},
		{Name: "PreVerificationGas", Type: "uint256"},
		{Name: "MaxFeePerGas", Type: "uint256"},
		{Name: "MaxPriorityFeePerGas", Type: "uint256"},
		{Name: "PaymasterAndData", Type: "bytes"},
		{Name: "Signature", Type: "bytes"},
	})
	bytes32, _                    = abi.NewType("bytes32", "", nil)
	uint256, _                    = abi.NewType("uint256", "", nil)
	bytes, _                      = abi.NewType("bytes", "", nil)
	validatePaymasterUserOpMethod = abi.NewMethod(
		"validatePaymasterUserOp",
		"validatePaymasterUserOp",
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{Name: "userOp", Type: userOp},
			{Name: "userOpHash", Type: bytes32},
			{Name: "maxCost", Type: uint256},
		},
		abi.Arguments{
			{Name: "context", Type: bytes},
			{Name: "deadline", Type: uint256},
		},
	)
	validatePaymasterUserOpSelector = hexutil.Encode(validatePaymasterUserOpMethod.ID)
)

type validatePaymasterUserOpOutput struct {
	Context []byte
}

func decodeValidatePaymasterUserOpOutput(ret any) (*validatePaymasterUserOpOutput, error) {
	hex, ok := ret.(string)
	if !ok {
		return nil, errors.New("validatePaymasterUserOp: cannot assert type: hex is not of type string")
	}
	data, err := hexutil.Decode(hex)
	if err != nil {
		return nil, fmt.Errorf("validatePaymasterUserOp: %s", err)
	}

	args, err := validatePaymasterUserOpMethod.Outputs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("validatePaymasterUserOp: %s", err)
	}
	if len(args) != 2 {
		return nil, fmt.Errorf(
			"validatePaymasterUserOp: invalid args length: expected 2, got %d",
			len(args),
		)
	}

	ctx, ok := args[0].([]byte)
	if !ok {
		return nil, errors.New("validatePaymasterUserOp: cannot assert type: hex is not of type string")
	}

	return &validatePaymasterUserOpOutput{
		Context: ctx,
	}, nil
}
