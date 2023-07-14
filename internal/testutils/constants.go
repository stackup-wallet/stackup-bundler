package testutils

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
)

var (
	OneETH                  = big.NewInt(1000000000000000000)
	DefaultUnstakeDelaySec  = uint32(86400)
	ValidAddress1           = common.HexToAddress("0x7357b8a705328FC283dF72D7Ac546895B596DC12")
	ValidAddress2           = common.HexToAddress("0x7357c9504B8686c008CCcD6ea47f1c21B7475dE3")
	ValidAddress3           = common.HexToAddress("0x7357C8D931e8cde8ea1b777Cf8578f4A7071f100")
	ChainID                 = big.NewInt(1)
	MaxOpsForUnstakedSender = 1
	StakedDepositInfo       = &entrypoint.IStakeManagerDepositInfo{
		Deposit:         big.NewInt(OneETH.Int64()),
		Staked:          true,
		Stake:           big.NewInt(OneETH.Int64()),
		UnstakeDelaySec: DefaultUnstakeDelaySec,
		WithdrawTime:    big.NewInt(time.Now().Unix()),
	}
	StakedZeroDepositInfo = &entrypoint.IStakeManagerDepositInfo{
		Deposit:         big.NewInt(0),
		Staked:          true,
		Stake:           big.NewInt(OneETH.Int64()),
		UnstakeDelaySec: DefaultUnstakeDelaySec,
		WithdrawTime:    big.NewInt(time.Now().Unix()),
	}
	NonStakedDepositInfo = &entrypoint.IStakeManagerDepositInfo{
		Deposit:         big.NewInt(OneETH.Int64()),
		Staked:          false,
		Stake:           big.NewInt(0),
		UnstakeDelaySec: uint32(0),
		WithdrawTime:    big.NewInt(0),
	}
	NonStakedZeroDepositInfo = &entrypoint.IStakeManagerDepositInfo{
		Deposit:         big.NewInt(0),
		Staked:          false,
		Stake:           big.NewInt(0),
		UnstakeDelaySec: uint32(0),
		WithdrawTime:    big.NewInt(0),
	}
)
