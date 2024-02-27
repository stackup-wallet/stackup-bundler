package state

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type OverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   *hexutil.Big                 `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

// OverrideSet is a set of accounts with customized state that can be applied during eth_call or
// debug_traceCall.
type OverrideSet map[common.Address]OverrideAccount

// Copy creates a deep copy of the OverrideSet.
func (os OverrideSet) Copy() (OverrideSet, error) {
	cpy := OverrideSet{}
	for k, v := range os {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		oa := OverrideAccount{}
		json.Unmarshal(b, &oa)
		cpy[k] = oa
	}

	return cpy, nil
}

// ParseOverrideData decodes a map into an OverrideSet and validates all the fields are correctly typed.
func ParseOverrideData(data map[string]any) (OverrideSet, error) {
	os := OverrideSet{}
	for key, value := range data {
		if !common.IsHexAddress(key) {
			return nil, fmt.Errorf("%w: %s", ErrBadKey, key)
		}

		b, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("%w %s: %w", ErrBadValue, key, err)
		}

		oa := OverrideAccount{}
		if err := json.Unmarshal(b, &oa); err != nil {
			return nil, fmt.Errorf("%w %s: %w", ErrBadValue, key, err)
		}

		os[common.HexToAddress(key)] = oa
	}
	return os, nil
}
