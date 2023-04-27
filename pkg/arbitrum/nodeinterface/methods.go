package nodeinterface

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	addressT, _ = abi.NewType("address", "", nil)
	boolT, _    = abi.NewType("bool", "", nil)
	bytesT, _   = abi.NewType("bytes", "", nil)
	uint64T, _  = abi.NewType("uint64", "", nil)
	uint256T, _ = abi.NewType("uint256", "", nil)

	GasEstimateL1ComponentMethod = abi.NewMethod(
		"gasEstimateL1Component",
		"gasEstimateL1Component",
		abi.Function,
		"",
		false,
		true,
		abi.Arguments{
			{Name: "to", Type: addressT},
			{Name: "contractCreation", Type: boolT},
			{Name: "data", Type: bytesT},
		},
		abi.Arguments{
			{Name: "gasEstimateForL1", Type: uint64T},
			{Name: "baseFee", Type: uint256T},
			{Name: "l1BaseFeeEstimate", Type: uint256T},
		},
	)
)

type GasEstimateL1ComponentOutput struct {
	GasEstimateForL1  uint64
	BaseFee           *big.Int
	L1BaseFeeEstimate *big.Int
}

func DecodeGasEstimateL1ComponentOutput(out any) (*GasEstimateL1ComponentOutput, error) {
	hex, ok := out.(string)
	if !ok {
		return nil, errors.New("gasEstimateL1Component: cannot assert type: hex is not of type string")
	}
	data, err := hexutil.Decode(hex)
	if err != nil {
		return nil, fmt.Errorf("gasEstimateL1Component: %s", err)
	}

	args, err := GasEstimateL1ComponentMethod.Outputs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("gasEstimateL1Component: %s", err)
	}

	return &GasEstimateL1ComponentOutput{
		GasEstimateForL1:  args[0].(uint64),
		BaseFee:           args[1].(*big.Int),
		L1BaseFeeEstimate: args[2].(*big.Int),
	}, nil
}
