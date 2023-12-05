package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// WithMaxBalanceOverride takes a set and appends an override for the given address to have a balance equal to
// max uint96.
func WithMaxBalanceOverride(acc common.Address, os OverrideSet) OverrideSet {
	if os == nil {
		os = OverrideSet{}
	}
	if _, ok := os[acc]; ok {
		return os
	}

	bal := hexutil.Big(*maxUint96)
	os[acc] = OverrideAccount{
		Balance: &bal,
	}

	return os
}
