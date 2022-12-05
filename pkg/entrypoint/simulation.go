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
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

var (
	// A dummy private key used to build *bind.TransactOpts for simulation.
	dummyPk, _ = crypto.GenerateKey()

	// A marker to delimit between account and paymaster simulation.
	markerOpCode = "NUMBER"

	// All opcodes executed at this depth are from the EntryPoint and allowed.
	allowedDepth = float64(1)

	// The gas opcode is only allowed if followed immediately by callOpcodes.
	gasOpCode = "GAS"

	// Only one create2 opcode is allowed if these two conditions are met:
	// 	1. op.initcode.length != 0
	// 	2. During account simulation (i.e. before markerOpCode)
	// create2OpCode = "CREATE2"

	// List of opcodes related to CALL.
	callOpcodes = mapset.NewSet(
		"CALL",
		"DELEGATECALL",
		"CALLCODE",
		"STATICCALL",
	)

	// List of opcodes not allowed during simulation for depth > allowedDepth (i.e. account, paymaster, or
	// contracts called by them).
	baseForbiddenOpCodes = mapset.NewSet(
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
		"CREATE",
		"COINBASE",
	)
)

// SimulateValidation makes a static call to Entrypoint.simulateValidation(userop) and returns the
// results without any state changes.
func SimulateValidation(rpc *rpc.Client, entryPoint common.Address, op *userop.UserOperation) (*SimulationResultRevert, error) {
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

	sim, simErr := newSimulationResultRevert(err)
	if simErr != nil {
		fo, foErr := newFailedOpRevert(err)
		if foErr != nil {
			return nil, fmt.Errorf("%s, %s", simErr, foErr)
		}
		return nil, errors.New(fo.Reason)
	}

	return sim, nil
}

type structLog struct {
	Depth   float64  `json:"depth"`
	Gas     float64  `json:"gas"`
	GasCost float64  `json:"gasCost"`
	Op      string   `json:"op"`
	Pc      float64  `json:"pc"`
	Stack   []string `json:"stack"`
}

type traceCallRes struct {
	Failed      bool        `json:"failed"`
	Gas         float64     `json:"gas"`
	ReturnValue []byte      `json:"returnValue"`
	StructLogs  []structLog `json:"structLogs"`
}

type traceCallReq struct {
	From common.Address `json:"from"`
	To   common.Address `json:"to"`
	Data hexutil.Bytes  `json:"data"`
}

type traceCallOpts struct {
	DisableStorage bool `json:"disableStorage"`
	DisableMemory  bool `json:"disableMemory"`
}

// TraceSimulateValidation makes a debug_traceCall to Entrypoint.simulateValidation(userop) and returns the
// results without any state changes.
func TraceSimulateValidation(rpc *rpc.Client, entryPoint common.Address, op *userop.UserOperation, chainID *big.Int) error {
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

	var res traceCallRes
	req := traceCallReq{
		From: common.HexToAddress("0x"),
		To:   entryPoint,
		Data: tx.Data(),
	}
	opts := traceCallOpts{
		DisableStorage: false,
		DisableMemory:  false,
	}
	if err := rpc.CallContext(context.Background(), &res, "debug_traceCall", &req, "latest", &opts); err != nil {
		return err
	}

	var prev structLog
	simFor := "account"
	for _, sl := range res.StructLogs {
		if sl.Depth == allowedDepth {
			if sl.Op == markerOpCode {
				simFor = "paymaster"
			}
			continue
		}

		if prev.Op == gasOpCode && !callOpcodes.Contains(sl.Op) {
			return fmt.Errorf("%s: uses opcode %s incorrectly", simFor, gasOpCode)
		}

		if baseForbiddenOpCodes.Contains(sl.Op) {
			return fmt.Errorf("%s: uses forbidden opcode %s", simFor, sl.Op)
		}

		prev = sl
	}

	return nil
}
