package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type TraceCallReq struct {
	From         common.Address `json:"from"`
	To           common.Address `json:"to"`
	Data         hexutil.Bytes  `json:"data"`
	MaxFeePerGas *hexutil.Big   `json:"maxFeePerGas,omitempty"`
}

type TraceCallOpts struct {
	Tracer string `json:"tracer"`
}

var (
	// A dummy private key used to build *bind.TransactOpts for simulation.
	DummyPk, _ = crypto.GenerateKey()
)
