package methods

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

var (
	ValidatePaymasterUserOpMethod = abi.NewMethod(
		"validatePaymasterUserOp",
		"validatePaymasterUserOp",
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{Name: "userOp", Type: userop.UserOpType},
			{Name: "userOpHash", Type: bytes32},
			{Name: "maxCost", Type: uint256},
		},
		abi.Arguments{
			{Name: "context", Type: bytes},
			{Name: "deadline", Type: uint256},
		},
	)
	ValidatePaymasterUserOpSelector = hexutil.Encode(ValidatePaymasterUserOpMethod.ID)
)

type validatePaymasterUserOpOutput struct {
	Context []byte
}

func DecodeValidatePaymasterUserOpOutput(ret any) (*validatePaymasterUserOpOutput, error) {
	hex, ok := ret.(string)
	if !ok {
		return nil, errors.New("validatePaymasterUserOp: cannot assert type: hex is not of type string")
	}
	data, err := hexutil.Decode(hex)
	if err != nil {
		return nil, fmt.Errorf("validatePaymasterUserOp: %s", err)
	}

	args, err := ValidatePaymasterUserOpMethod.Outputs.Unpack(data)
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
