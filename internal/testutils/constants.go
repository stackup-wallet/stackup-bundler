package testutils

import (
	entrypointV06 "github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/bindings/v06"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
)

var (
	OneETH                  = big.NewInt(1000000000000000000)
	DefaultUnstakeDelaySec  = uint32(86400)
	ValidAddress1           = common.HexToAddress("0x7357b8a705328FC283dF72D7Ac546895B596DC12")
	ValidAddress2           = common.HexToAddress("0x7357c9504B8686c008CCcD6ea47f1c21B7475dE3")
	ValidAddress3           = common.HexToAddress("0x7357C8D931e8cde8ea1b777Cf8578f4A7071f100")
	ValidAddress4           = common.HexToAddress("0x73574a159D05d20FF50D5504057D5C86f2d02a45")
	ValidAddress5           = common.HexToAddress("0x7357C1Fc72a14399cb845f2f71421B4CE7eCE608")
	ChainID                 = big.NewInt(1)
	MaxOpsForUnstakedSender = 1
	StakedDepositInfo       = &entrypointV06.IStakeManagerDepositInfo{
		Deposit:         big.NewInt(OneETH.Int64()),
		Staked:          true,
		Stake:           big.NewInt(OneETH.Int64()),
		UnstakeDelaySec: DefaultUnstakeDelaySec,
		WithdrawTime:    big.NewInt(time.Now().Unix()),
	}
	StakedZeroDepositInfo = &entrypointV06.IStakeManagerDepositInfo{
		Deposit:         big.NewInt(0),
		Staked:          true,
		Stake:           big.NewInt(OneETH.Int64()),
		UnstakeDelaySec: DefaultUnstakeDelaySec,
		WithdrawTime:    big.NewInt(time.Now().Unix()),
	}
	NonStakedDepositInfo = &entrypointV06.IStakeManagerDepositInfo{
		Deposit:         big.NewInt(OneETH.Int64()),
		Staked:          false,
		Stake:           big.NewInt(0),
		UnstakeDelaySec: uint32(0),
		WithdrawTime:    big.NewInt(0),
	}
	NonStakedZeroDepositInfo = &entrypointV06.IStakeManagerDepositInfo{
		Deposit:         big.NewInt(0),
		Staked:          false,
		Stake:           big.NewInt(0),
		UnstakeDelaySec: uint32(0),
		WithdrawTime:    big.NewInt(0),
	}

	pk, _       = crypto.GenerateKey()
	DummyEOA, _ = signer.New(hexutil.Encode(crypto.FromECDSA(pk))[2:])
	MockHash    = "0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddead"
)
