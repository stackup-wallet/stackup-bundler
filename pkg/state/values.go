package state

import (
	"errors"
	"math/big"
)

var (
	maxUint96, _ = big.NewInt(0).SetString("79228162514264337593543950335", 10)

	ErrBadKey   = errors.New("cannot decode key to address")
	ErrBadValue = errors.New("cannot decode override account")
)
