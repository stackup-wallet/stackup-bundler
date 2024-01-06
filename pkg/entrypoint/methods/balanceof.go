package methods

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	BalanceOfMethod = abi.NewMethod(
		"balanceOf",
		"balanceOf",
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{Name: "account", Type: address},
		},
		abi.Arguments{
			{Type: uint256},
		},
	)
	BalanceOfSelector = hexutil.Encode(BalanceOfMethod.ID)
)
