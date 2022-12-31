package entrypoint

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

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
	)
)

// SimulateValidation makes a static call to Entrypoint.simulateValidation(userop) and returns the
// results without any state changes.
func SimulateValidation(
	rpc *rpc.Client,
	entryPoint common.Address,
	op *userop.UserOperation,
) (*ValidationResultRevert, error) {
	ep, err := NewEntrypoint(entryPoint, ethclient.NewClient(rpc))
	if err != nil {
		return nil, err
	}

	var res []interface{}
	rawCaller := &EntrypointRaw{Contract: ep}
	err = rawCaller.Call(nil, &res, "simulateValidation", UserOperation(*op))
	if err == nil {
		return nil, errors.New("unexpected result from simulateValidation")
	}

	sim, simErr := newValidationResultRevert(err)
	if simErr != nil {
		fo, foErr := newFailedOpRevert(err)
		if foErr != nil {
			return nil, fmt.Errorf("%s, %s", simErr, foErr)
		}
		return nil, errors.New(fo.Reason)
	}

	return sim, nil
}

type traceCallReq struct {
	From common.Address `json:"from"`
	To   common.Address `json:"to"`
	Data hexutil.Bytes  `json:"data"`
}

type traceCallOpts struct {
	Tracer string `json:"tracer"`
}

// TraceSimulateValidation makes a debug_traceCall to Entrypoint.simulateValidation(userop) and returns the
// results without any state changes.
func TraceSimulateValidation(
	rpc *rpc.Client,
	entryPoint common.Address,
	op *userop.UserOperation,
	chainID *big.Int,
	customTracer string,
) error {
	ep, err := NewEntrypoint(entryPoint, ethclient.NewClient(rpc))
	if err != nil {
		return err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(dummyPk, chainID)
	if err != nil {
		return err
	}
	auth.GasLimit = math.MaxUint64
	auth.NoSend = true
	tx, err := ep.SimulateValidation(auth, UserOperation(*op))
	if err != nil {
		return err
	}

	req := traceCallReq{
		From: common.HexToAddress("0x"),
		To:   entryPoint,
		Data: tx.Data(),
	}

	var res tracer.BundlerCollectorReturn
	opts := traceCallOpts{
		Tracer: customTracer,
	}
	if err := rpc.CallContext(context.Background(), &res, "debug_traceCall", &req, "latest", &opts); err != nil {
		return err
	}

	var factoryOpCodes, accountOpCodes, paymasterOpCodes tracer.Counts
	if len(res.NumberLevels) != 3 {
		return fmt.Errorf("unexpected tracing result for op: %s", op.GetUserOpHash(entryPoint, chainID))
	}
	factoryOpCodes = res.NumberLevels[factoryNumberLevel].Opcodes
	accountOpCodes = res.NumberLevels[accountNumberLevel].Opcodes
	paymasterOpCodes = res.NumberLevels[paymasterNumberLevel].Opcodes

	for opcode := range factoryOpCodes {
		if bannedOpCodes.Contains(opcode) {
			return fmt.Errorf("factory contains banned opcode: %s", opcode)
		}
	}

	for opcode := range accountOpCodes {
		if bannedOpCodes.Contains(opcode) {
			return fmt.Errorf("account contains banned opcode: %s", opcode)
		}
	}

	for opcode := range paymasterOpCodes {
		if bannedOpCodes.Contains(opcode) {
			return fmt.Errorf("paymaster contains banned opcode: %s", opcode)
		}
	}

	create2Count, ok := factoryOpCodes[create2OpCode]
	if ok && (create2Count > 1 || len(op.InitCode) == 0) {
		return fmt.Errorf("factory with too many %s", create2OpCode)
	}

	_, ok = accountOpCodes[create2OpCode]
	if ok {
		return fmt.Errorf("account uses banned %s opcode: %s", create2OpCode, op.Sender.String())
	}

	_, ok = paymasterOpCodes[create2OpCode]
	if ok {
		return fmt.Errorf("paymaster uses banned %s opcode: %s", create2OpCode, op.GetPaymaster())
	}

	return nil
}
