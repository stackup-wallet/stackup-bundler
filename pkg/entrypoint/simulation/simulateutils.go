package simulation

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
)

// EntityStakes provides a mapping for encountered entity addresses and their stake info on the EntryPoint.
type EntityStakes map[common.Address]*entrypoint.IStakeManagerDepositInfo

var (
	// Only one create2 opcode is allowed if these two conditions are met:
	// 	1. op.initcode.length != 0
	// 	2. During account simulation (i.e. before markerOpCode)
	create2OpCode = "CREATE2"

	// List of opcodes not allowed during simulation for depth > 1 (i.e. account, paymaster, or contracts
	// called by them).
	bannedOpCodes = mapset.NewSet(
		"GASPRICE",
		"GASLIMIT",
		"DIFFICULTY",
		"TIMESTAMP",
		"BASEFEE",
		"BLOCKHASH",
		"NUMBER",
		"ORIGIN",
		"GAS",
		"CREATE",
		"COINBASE",
		"SELFDESTRUCT",
	)

	// List of opcodes not allowed during validation for unstaked entities.
	bannedUnstakedOpCodes = mapset.NewSet(
		"SELFBALANCE",
		"BALANCE",
	)

	revertOpCode = "REVERT"
	returnOpCode = "RETURN"

	// Precompiled contract that performs secp256r1 signature verification. See
	// https://github.com/ethereum/RIPs/blob/master/RIPS/rip-7212.md
	rip7212precompile = common.HexToAddress("0x100")
)
