package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	maxUint96, _ = big.NewInt(0).SetString("79228162514264337593543950335", 10)
)

// WithZeroAddressOverride takes a set and appends an override for the zero address to have a balance equal to
// max uint96.
func WithZeroAddressOverride(os OverrideSet) OverrideSet {
	if os == nil {
		os = OverrideSet{}
	}

	bal := hexutil.Big(*maxUint96)
	os[common.HexToAddress("0x")] = OverrideAccount{
		Balance: &bal,
	}

	return os
}
