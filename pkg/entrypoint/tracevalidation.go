package entrypoint

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// TraceSimulateValidation makes a debug_traceCall to Entrypoint.simulateValidation(userop) and returns the
// results without any state changes.
func TraceSimulateValidation(
	rpc *rpc.Client,
	entryPoint common.Address,
	op *userop.UserOperation,
	chainID *big.Int,
	customTracer string,
	stakes EntityStakes,
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

	var res tracer.BundlerCollectorReturn
	req := traceCallReq{
		From: common.HexToAddress("0x"),
		To:   entryPoint,
		Data: tx.Data(),
	}
	opts := traceCallOpts{
		Tracer: customTracer,
	}
	if err := rpc.CallContext(context.Background(), &res, "debug_traceCall", &req, "latest", &opts); err != nil {
		return err
	}

	entityMap, err := newEntityMap(op, &res, stakes)
	if err != nil {
		return err
	}

	for title, entity := range entityMap {
		for opcode := range entity.Info.Opcodes {
			if bannedOpCodes.Contains(opcode) {
				return fmt.Errorf("%s uses banned opcode: %s", title, opcode)
			}
		}
	}

	create2Count, ok := entityMap["factory"].Info.Opcodes[create2OpCode]
	if ok && (create2Count > 1 || len(op.InitCode) == 0) {
		return fmt.Errorf("factory with too many %s", create2OpCode)
	}
	_, ok = entityMap["account"].Info.Opcodes[create2OpCode]
	if ok {
		return fmt.Errorf("account uses banned opcode: %s", create2OpCode)
	}
	_, ok = entityMap["paymaster"].Info.Opcodes[create2OpCode]
	if ok {
		return fmt.Errorf("paymaster uses banned opcode: %s", create2OpCode)
	}

	slots := parseEntitySlots(stakes, res.Keccak)
	for title, entity := range entityMap {
		if err := validateEntityStorage(
			title,
			op,
			entryPoint,
			slots,
			entity.Info.Access,
			entity.Address,
			entity.IsStaked,
		); err != nil {
			return err
		}
	}

	callStack := parseCallStack(res.Calls)
	for _, call := range callStack {
		if call.Method == validatePaymasterUserOpSelector {
			out, err := decodeValidatePaymasterUserOpOutput(call.Return)
			if err != nil {
				return fmt.Errorf(
					"unexpected tracing result for op: %s, %s",
					op.GetUserOpHash(entryPoint, chainID),
					err,
				)
			}

			if len(out.Context) != 0 && !entityMap["paymaster"].IsStaked {
				return errors.New("unstaked paymaster must not return context")
			}
		}
	}

	return nil
}
