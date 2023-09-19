package simulation

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
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/methods"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/utils"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// TraceSimulateValidation makes a debug_traceCall to Entrypoint.simulateValidation(userop) and returns an
// array of all the interacted contracts touched by entities during the trace.
func TraceSimulateValidation(
	rpc *rpc.Client,
	entryPoint common.Address,
	op *userop.UserOperation,
	chainID *big.Int,
	stakes EntityStakes,
) ([]common.Address, error) {
	ep, err := entrypoint.NewEntrypoint(entryPoint, ethclient.NewClient(rpc))
	if err != nil {
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(utils.DummyPk, chainID)
	if err != nil {
		return nil, err
	}
	auth.GasLimit = math.MaxUint64
	auth.NoSend = true
	tx, err := ep.SimulateValidation(auth, entrypoint.UserOperation(*op))
	if err != nil {
		return nil, err
	}

	var res tracer.BundlerCollectorReturn
	req := utils.TraceCallReq{
		From:         common.HexToAddress("0x"),
		To:           entryPoint,
		Data:         tx.Data(),
		MaxFeePerGas: hexutil.Big(*op.MaxFeePerGas),
	}
	opts := utils.TraceCallOpts{
		Tracer:         tracer.Loaded.BundlerCollectorTracer,
		StateOverrides: utils.DefaultStateOverrides,
	}
	if err := rpc.CallContext(context.Background(), &res, "debug_traceCall", &req, "latest", &opts); err != nil {
		return nil, err
	}

	knownEntity, err := newKnownEntity(op, &res, stakes)
	if err != nil {
		return nil, err
	}

	ic := mapset.NewSet[common.Address]()
	for title, entity := range knownEntity {
		for opcode := range entity.Info.Opcodes {
			if bannedOpCodes.Contains(opcode) {
				return nil, fmt.Errorf("%s uses banned opcode: %s", title, opcode)
			}
		}

		for addrHex := range entity.Info.ContractSize {
			ic.Add(common.HexToAddress(addrHex))
		}
	}

	create2Count, ok := knownEntity["factory"].Info.Opcodes[create2OpCode]
	if ok && (create2Count > 1 || len(op.InitCode) == 0) {
		return nil, fmt.Errorf("factory with too many %s", create2OpCode)
	}
	_, ok = knownEntity["account"].Info.Opcodes[create2OpCode]
	if ok {
		return nil, fmt.Errorf("account uses banned opcode: %s", create2OpCode)
	}
	_, ok = knownEntity["paymaster"].Info.Opcodes[create2OpCode]
	if ok {
		return nil, fmt.Errorf("paymaster uses banned opcode: %s", create2OpCode)
	}

	slotsByEntity := newStorageSlotsByEntity(stakes, res.Keccak)
	for title, entity := range knownEntity {
		v := &storageSlotsValidator{
			Op:              op,
			EntryPoint:      entryPoint,
			SenderSlots:     slotsByEntity[op.Sender],
			FactoryIsStaked: knownEntity["factory"].IsStaked,
			EntityName:      title,
			EntityAddr:      entity.Address,
			EntityAccess:    entity.Info.Access,
			EntitySlots:     slotsByEntity[entity.Address],
			EntityIsStaked:  entity.IsStaked,
		}
		if err := v.Process(); err != nil {
			return nil, err
		}
	}

	callStack := newCallStack(res.Calls)
	for _, call := range callStack {
		if call.Method == methods.ValidatePaymasterUserOpSelector {
			out, err := methods.DecodeValidatePaymasterUserOpOutput(call.Return)
			if err != nil {
				return nil, fmt.Errorf(
					"unexpected tracing result for op: %s, %s",
					op.GetUserOpHash(entryPoint, chainID),
					err,
				)
			}

			if len(out.Context) != 0 && !knownEntity["paymaster"].IsStaked {
				return nil, errors.New("unstaked paymaster must not return context")
			}
		}
	}

	return ic.ToSlice(), nil
}
