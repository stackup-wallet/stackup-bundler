package testutils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
)

func MockGetCode(addr common.Address) ([]byte, error) {
	return MockByteCode, nil
}

func MockGetCodeZero(addr common.Address) ([]byte, error) {
	return []byte{}, nil
}

func MockGetStake(addr common.Address) (*entrypoint.IStakeManagerDepositInfo, error) {
	return StakedDepositInfo, nil
}

func MockGetStakeZeroDeposit(addr common.Address) (*entrypoint.IStakeManagerDepositInfo, error) {
	return StakedZeroDepositInfo, nil
}

func MockGetNotStake(addr common.Address) (*entrypoint.IStakeManagerDepositInfo, error) {
	return NonStakedDepositInfo, nil
}

func MockGetNotStakeZeroDeposit(addr common.Address) (*entrypoint.IStakeManagerDepositInfo, error) {
	return NonStakedZeroDepositInfo, nil
}

func GetMockBaseFeeFunc(val *big.Int) func() (*big.Int, error) {
	return func() (*big.Int, error) {
		return val, nil
	}
}
