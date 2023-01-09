package entrypoint

import (
	"fmt"
	"math/big"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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

type EntityMap map[string]struct {
	Address  common.Address
	Info     tracer.NumberLevelInfo
	IsStaked bool
}

func newEntityMap(
	op *userop.UserOperation,
	res *tracer.BundlerCollectorReturn,
	stakes EntityStakes,
) (EntityMap, error) {
	if len(res.NumberLevels) != 3 {
		return nil, fmt.Errorf("unexpected NumberLevels length in tracing result: %d", len(res.NumberLevels))
	}

	return EntityMap{
		"factory": {
			Address:  op.GetFactory(),
			Info:     res.NumberLevels[factoryNumberLevel],
			IsStaked: stakes[op.GetFactory()] != nil && stakes[op.GetFactory()].Staked,
		},
		"account": {
			Address:  op.Sender,
			Info:     res.NumberLevels[accountNumberLevel],
			IsStaked: stakes[op.Sender] != nil && stakes[op.Sender].Staked,
		},
		"paymaster": {
			Address:  op.GetPaymaster(),
			Info:     res.NumberLevels[paymasterNumberLevel],
			IsStaked: stakes[op.GetPaymaster()] != nil && stakes[op.GetPaymaster()].Staked,
		},
	}, nil
}

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
			if _, ok := entitySlots[addr]; !ok {
				entitySlots[addr] = mapset.NewSet[string]()
			}

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

func associatedWith(slots EntitySlots, addr common.Address, slot string) bool {
	entitySlot, entitySlotOk := slots[addr]
	if !entitySlotOk {
		return false
	}

	slotN, _ := big.NewInt(0).SetString(fmt.Sprintf("0x%s", slot), 0)
	for _, k := range entitySlot.ToSlice() {
		kn, _ := big.NewInt(0).SetString(fmt.Sprintf("0x%s", k), 0)
		if slotN.Cmp(kn) >= 0 && slotN.Cmp(big.NewInt(0).Add(kn, big.NewInt(128))) <= 0 {
			return true
		}
	}

	return false
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
		accessTypes := map[string]tracer.Counts{
			"read":  access.Reads,
			"write": access.Writes,
		}
		for key, slotCount := range accessTypes {
			for slot := range slotCount {
				if associatedWith(slots, op.Sender, slot) {
					if len(op.InitCode) > 0 {
						mustStakeSlot = slot
					} else {
						continue
					}
				} else if associatedWith(slots, entityAddr, slot) || addr == entityAddr {
					mustStakeSlot = slot
				} else {
					return fmt.Errorf("%s has forbidden %s to %s slot %s", entityName, key, nameAddr(op, addr), slot)
				}
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
