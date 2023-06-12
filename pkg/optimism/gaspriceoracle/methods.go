package gaspriceoracle

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	bytesT, _   = abi.NewType("bytes", "", nil)
	uint256T, _ = abi.NewType("uint256", "", nil)

	GetL1FeeMethod = abi.NewMethod(
		"getL1Fee",
		"getL1Fee",
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{Name: "data", Type: bytesT},
		},
		abi.Arguments{
			{Name: "fee", Type: uint256T},
		},
	)
)

func DecodeGetL1FeeMethodOutput(out any) (*big.Int, error) {
	hex, ok := out.(string)
	if !ok {
		return nil, errors.New("getL1Fee: cannot assert type: hex is not of type string")
	}
	data, err := hexutil.Decode(hex)
	if err != nil {
		return nil, fmt.Errorf("getL1Fee: %s", err)
	}

	args, err := GetL1FeeMethod.Outputs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("getL1Fee: %s", err)
	}

	return args[0].(*big.Int), nil
}
