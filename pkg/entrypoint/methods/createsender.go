package methods

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	CreateSenderMethod = abi.NewMethod(
		"createSender",
		"createSender",
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{Name: "initCode", Type: bytes},
		},
		abi.Arguments{
			{Name: "sender", Type: address},
		},
	)
	CreateSenderSelector = hexutil.Encode(CreateSenderMethod.ID)
)
