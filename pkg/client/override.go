package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Named UserOperation type for jsonrpc package.
type userOperation map[string]any

// Named StateOverride type for jsonrpc package.
type optional_stateOverride map[string]any

type OverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

// StateOverride is a set of accounts with customized state that can be applied during gas estimations.
type StateOverride map[common.Address]OverrideAccount
