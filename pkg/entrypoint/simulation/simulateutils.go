package simulation

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
)

// EntityStakes provides a mapping for encountered entity addresses and their stake info on the EntryPoint.
type EntityStakes map[common.Address]*entrypoint.IStakeManagerDepositInfo

type traceCallReq struct {
	From common.Address `json:"from"`
	To   common.Address `json:"to"`
	Data hexutil.Bytes  `json:"data"`
}

type traceCallOpts struct {
	Tracer string `json:"tracer"`
}

var (
	// A dummy private key used to build *bind.TransactOpts for simulation.
	dummyPk, _ = crypto.GenerateKey()

	// Up to the first number marker represents factory validation.
	factoryNumberLevel = 0

	// After the first number marker and before the second represents account validation.
	accountNumberLevel = 1

	// After the second number marker represents paymaster validation.
	paymasterNumberLevel = 2

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
		"SELFBALANCE",
		"BALANCE",
		"ORIGIN",
		"GAS",
		"CREATE",
		"COINBASE",
		"SELFDESTRUCT",
	)

	revertOpCode = "REVERT"
	returnOpCode = "RETURN"
)
