package gaspriceoracle

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	uint256T, _ = abi.NewType("uint256", "", nil)

	GetL1FeeUpperBoundMethod = abi.NewMethod(
		"getL1FeeUpperBound",
		"getL1FeeUpperBound",
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{Name: "data", Type: uint256T},
		},
		abi.Arguments{
			{Name: "fee", Type: uint256T},
		},
	)
)

func DecodeGetL1FeeUpperBoundOutput(out any) (*big.Int, error) {
	hex, ok := out.(string)
	if !ok {
		return nil, errors.New("getL1Fee: cannot assert type: hex is not of type string")
	}
	data, err := hexutil.Decode(hex)
	if err != nil {
		return nil, fmt.Errorf("getL1Fee: %s", err)
	}

	args, err := GetL1FeeUpperBoundMethod.Outputs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("getL1Fee: %s", err)
	}

	return args[0].(*big.Int), nil
}
