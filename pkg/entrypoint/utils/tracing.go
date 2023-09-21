package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type TraceCallReq struct {
	From         common.Address `json:"from"`
	To           common.Address `json:"to"`
	Data         hexutil.Bytes  `json:"data"`
	MaxFeePerGas hexutil.Big    `json:"maxFeePerGas"`
}

type TraceStateOverrides struct {
	Balance hexutil.Big `json:"balance"`
}

type TraceCallOpts struct {
	Tracer         string                         `json:"tracer"`
	StateOverrides map[string]TraceStateOverrides `json:"stateOverrides"`
}

var (
	// A dummy private key used to build *bind.TransactOpts for simulation.
	DummyPk, _ = crypto.GenerateKey()

	maxUint96, _ = big.NewInt(0).SetString("79228162514264337593543950335", 10)

	// A default state override to ensure the zero address always has sufficient funds.
	DefaultStateOverrides = map[string]TraceStateOverrides{
		common.HexToAddress("0x").Hex(): {Balance: hexutil.Big(*maxUint96)},
	}
)
