package entrypoint

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/internal/utils"
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

	revertOpCode = "REVERT"
	returnOpCode = "RETURN"
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

type callEntry struct {
	To     common.Address
	Type   string
	Method string
	Revert any
	Return any
	Value  *big.Int
}

type EntityStakes = map[common.Address]*IStakeManagerDepositInfo

type EntitySlots map[common.Address]mapset.Set[string]

func parseCallStack(calls []tracer.CallInfo) []*callEntry {
	out := []*callEntry{}
	stack := utils.NewStack[tracer.CallInfo]()
	for _, call := range calls {
		if call.Type == revertOpCode || call.Type == returnOpCode {
			top, _ := stack.Pop()

			if strings.Contains(top.Type, "CREATE") {
				// TODO: implement...
			} else if call.Type == revertOpCode {
				// TODO: implement...
			} else {
				out = append(out, &callEntry{
					To:     top.To,
					Type:   top.Type,
					Method: top.Method,
					Return: call.Data,
				})
			}
		} else {
			stack.Push(call)
		}
	}

	return out
}

func parseEntitySlots(stakes EntityStakes, keccak []string) EntitySlots {
	entitySlots := make(EntitySlots)

	for _, k := range keccak {
		value := common.Bytes2Hex(crypto.Keccak256(common.Hex2Bytes(k[2:])))

		for addr := range stakes {
			if addr == common.HexToAddress("0x") {
				continue
			}

			entitySlots[addr] = mapset.NewSet[string]()
			addrPadded := hexutil.Encode(common.LeftPadBytes(addr.Bytes(), 32))
			if strings.HasPrefix(k, addrPadded) {
				entitySlots[addr].Add(value)
			}
		}
	}

	return entitySlots
}

func nameAddr(op *userop.UserOperation, addr common.Address) string {
	if addr == op.Sender {
		return "account"
	} else if addr == op.GetPaymaster() {
		return "paymaster"
	} else if addr == op.GetFactory() {
		return "factory"
	} else {
		return addr.String()
	}
}

func entityIsStaked(stakes EntityStakes, entity common.Address) bool {
	entityStake := stakes[entity]
	return entityStake != nil && entityStake.Staked
}

func validateEntityStorage(
	entityName string,
	op *userop.UserOperation,
	entryPoint common.Address,
	slots EntitySlots,
	entityAccess tracer.AccessMap,
	entityAddr common.Address,
	entityIsStaked bool,
) error {
	for addr, access := range entityAccess {
		if addr == op.Sender || addr == entryPoint {
			continue
		}

		var mustStakeSlot string
		for slot := range access.Writes {
			senderSlot, ok := slots[op.Sender]
			if ok && senderSlot.Contains(slot) && len(op.InitCode) == 0 {
				continue
			}

			entitySlot, ok := slots[entityAddr]
			if (ok && entitySlot.Contains(slot)) || addr == entityAddr {
				mustStakeSlot = slot
			} else {
				return fmt.Errorf("%s has forbidden write to %s slot %s", entityName, nameAddr(op, addr), slot)
			}
		}
		for slot := range access.Reads {
			senderSlot, ok := slots[op.Sender]
			if ok && senderSlot.Contains(slot) && len(op.InitCode) == 0 {
				continue
			}

			entitySlot, ok := slots[entityAddr]
			if (ok && entitySlot.Contains(slot)) || addr == entityAddr {
				mustStakeSlot = slot
			} else {
				return fmt.Errorf("%s has forbidden read to %s slot %s", entityName, nameAddr(op, addr), slot)
			}
		}

		if mustStakeSlot != "" && !entityIsStaked {
			return fmt.Errorf(
				"unstaked %s accessed %s slot %s",
				entityName,
				nameAddr(op, addr),
				mustStakeSlot,
			)
		}
	}

	return nil
}

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

	if len(res.NumberLevels) != 3 {
		return fmt.Errorf("unexpected tracing result for op: %s", op.GetUserOpHash(entryPoint, chainID))
	}

	factoryOpCodes := res.NumberLevels[factoryNumberLevel].Opcodes
	accountOpCodes := res.NumberLevels[accountNumberLevel].Opcodes
	paymasterOpCodes := res.NumberLevels[paymasterNumberLevel].Opcodes

	for opcode := range factoryOpCodes {
		if bannedOpCodes.Contains(opcode) {
			return fmt.Errorf("factory uses banned opcode: %s", opcode)
		}
	}

	for opcode := range accountOpCodes {
		if bannedOpCodes.Contains(opcode) {
			return fmt.Errorf("account uses banned opcode: %s", opcode)
		}
	}

	for opcode := range paymasterOpCodes {
		if bannedOpCodes.Contains(opcode) {
			return fmt.Errorf("paymaster uses banned opcode: %s", opcode)
		}
	}

	create2Count, ok := factoryOpCodes[create2OpCode]
	if ok && (create2Count > 1 || len(op.InitCode) == 0) {
		return fmt.Errorf("factory with too many %s", create2OpCode)
	}

	_, ok = accountOpCodes[create2OpCode]
	if ok {
		return fmt.Errorf("account uses banned opcode: %s", create2OpCode)
	}

	_, ok = paymasterOpCodes[create2OpCode]
	if ok {
		return fmt.Errorf("paymaster uses banned opcode: %s", create2OpCode)
	}

	factory := op.GetFactory()
	paymaster := op.GetPaymaster()
	slots := parseEntitySlots(stakes, res.Keccak)
	if err := validateEntityStorage(
		"factory",
		op,
		entryPoint,
		slots,
		res.NumberLevels[factoryNumberLevel].Access,
		factory,
		entityIsStaked(stakes, factory),
	); err != nil {
		return err
	}
	if err := validateEntityStorage(
		"account",
		op,
		entryPoint,
		slots,
		res.NumberLevels[accountNumberLevel].Access,
		op.Sender,
		entityIsStaked(stakes, op.Sender),
	); err != nil {
		return err
	}
	if err := validateEntityStorage(
		"paymaster",
		op,
		entryPoint,
		slots,
		res.NumberLevels[paymasterNumberLevel].Access,
		paymaster,
		entityIsStaked(stakes, paymaster),
	); err != nil {
		return err
	}

	callStack := parseCallStack(res.Calls)
	for _, call := range callStack {
		if call.Method == validatePaymasterUserOpSelector {
			out, err := decodeValidatePaymasterUserOpOutput(call.Return)
			if err != nil {
				return fmt.Errorf("unexpected tracing result for op: %s", err)
			}
			if len(out.Context) != 0 && !entityIsStaked(stakes, paymaster) {
				return errors.New("unstaked paymaster must not return context")
			}
		}
	}

	return nil
}
