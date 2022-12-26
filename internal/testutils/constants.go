package testutils

import (
	"math/big"
	"time"

	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
)

var (
	OneETH                 = big.NewInt(1000000000000000000)
	DefaultUnstakeDelaySec = uint32(86400)
	StakedDepositInfo      = &entrypoint.IStakeManagerDepositInfo{
		Deposit:         big.NewInt(OneETH.Int64()),
		Staked:          true,
		Stake:           big.NewInt(OneETH.Int64()),
		UnstakeDelaySec: DefaultUnstakeDelaySec,
		WithdrawTime:    uint64(time.Now().Unix()),
	}
	StakedZeroDepositInfo = &entrypoint.IStakeManagerDepositInfo{
		Deposit:         big.NewInt(0),
		Staked:          true,
		Stake:           big.NewInt(OneETH.Int64()),
		UnstakeDelaySec: DefaultUnstakeDelaySec,
		WithdrawTime:    uint64(time.Now().Unix()),
	}
	NonStakedDepositInfo = &entrypoint.IStakeManagerDepositInfo{
		Deposit:         big.NewInt(OneETH.Int64()),
		Staked:          false,
		Stake:           big.NewInt(0),
		UnstakeDelaySec: uint32(0),
		WithdrawTime:    uint64(0),
	}
	NonStakedZeroDepositInfo = &entrypoint.IStakeManagerDepositInfo{
		Deposit:         big.NewInt(0),
		Staked:          false,
		Stake:           big.NewInt(0),
		UnstakeDelaySec: uint32(0),
		WithdrawTime:    uint64(0),
	}
)
