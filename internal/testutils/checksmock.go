package testutils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func MockGetCode(addr common.Address) ([]byte, error) {
	return MockByteCode, nil
}

func MockGetCodeZero(addr common.Address) ([]byte, error) {
	return []byte{}, nil
}

func GetMockBaseFeeFunc(val *big.Int) func() (*big.Int, error) {
	return func() (*big.Int, error) {
		return val, nil
	}
}
