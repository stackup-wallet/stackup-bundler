package simulation

import (
	"fmt"
	"math/big"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stackup-wallet/stackup-bundler/pkg/altmempools"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

var (
	accessModeRead       = "read"
	accessModeWrite      = "write"
	associatedSlotOffset = big.NewInt(128)
)

type storageSlots mapset.Set[string]

type storageSlotsByEntity map[common.Address]storageSlots

func newStorageSlotsByEntity(stakes EntityStakes, keccak []string) storageSlotsByEntity {
	storageSlotsByEntity := make(storageSlotsByEntity)

	for _, k := range keccak {
		value := hexutil.Encode(crypto.Keccak256(common.Hex2Bytes(k[2:])))

		for addr := range stakes {
			if addr == common.HexToAddress("0x") {
				continue
			}
			if _, ok := storageSlotsByEntity[addr]; !ok {
				storageSlotsByEntity[addr] = mapset.NewSet[string]()
			}

			addrPadded := hexutil.Encode(common.LeftPadBytes(addr.Bytes(), 32))
			if strings.HasPrefix(k, addrPadded) {
				storageSlotsByEntity[addr].Add(value)
			}
		}
	}

	return storageSlotsByEntity
}

type storageSlotsValidator struct {
	// Global parameters
	Op          *userop.UserOperation
	EntryPoint  common.Address
	AltMempools *altmempools.Directory

	// Parameters of specific entities required for all validation
	SenderSlots     storageSlots
	FactoryIsStaked bool

	// Parameters of the entity under validation
	EntityName            string
	EntityAddr            common.Address
	EntityAccessMap       tracer.AccessMap
	EntityContractSizeMap tracer.ContractSizeMap
	EntitySlots           storageSlots
	EntityIsStaked        bool
}

func isAssociatedWith(entitySlots storageSlots, slot string) bool {
	slotBN, _ := big.NewInt(0).SetString(slot, 0)
	for _, entitySlot := range entitySlots.ToSlice() {
		entitySlotBN, _ := big.NewInt(0).SetString(entitySlot, 0)
		maxAssocSlotBN := big.NewInt(0).Add(entitySlotBN, associatedSlotOffset)
		if slotBN.Cmp(entitySlotBN) >= 0 && slotBN.Cmp(maxAssocSlotBN) <= 0 {
			return true
		}
	}

	return false
}

func (v *storageSlotsValidator) Process() ([]string, error) {
	senderSlots := v.SenderSlots
	if senderSlots == nil {
		senderSlots = mapset.NewSet[string]()
	}
	entitySlots := v.EntitySlots
	if entitySlots == nil {
		entitySlots = mapset.NewSet[string]()
	}
	altMempoolIds := []string{}

	for ca, csi := range v.EntityContractSizeMap {
		if ca != v.Op.Sender && csi.ContractSize == 0 {
			return altMempoolIds, fmt.Errorf(
				"%s uses %s on an address with no deployed code: %s",
				v.EntityName,
				csi.Opcode,
				ca,
			)
		}
	}

	for addr, access := range v.EntityAccessMap {
		if addr == v.Op.Sender || addr == v.EntryPoint {
			continue
		}

		var mustStakeSlot string
		accessTypes := map[string]any{
			accessModeRead:  access.Reads,
			accessModeWrite: access.Writes,
		}
		for mode, val := range accessTypes {
			slots := []string{}
			if readMap, ok := val.(tracer.HexMap); ok {
				for slot := range readMap {
					slots = append(slots, slot)
				}
			} else if writeMap, ok := val.(tracer.Counts); ok {
				for slot := range writeMap {
					slots = append(slots, slot)
				}
			} else {
				return altMempoolIds, fmt.Errorf("cannot decode %s access type: %+v", mode, val)
			}

			for _, slot := range slots {
				if isAssociatedWith(senderSlots, slot) {
					if (len(v.Op.InitCode) > 0 && !v.FactoryIsStaked) ||
						(len(v.Op.InitCode) > 0 && v.FactoryIsStaked && v.EntityAddr != v.Op.Sender) {
						mustStakeSlot = slot
					} else {
						continue
					}
				} else if isAssociatedWith(entitySlots, slot) || mode == accessModeRead {
					mustStakeSlot = slot
				} else if ids := v.AltMempools.HasInvalidStorageAccessException(
					v.EntityName,
					addr2KnownEntity(v.Op, addr),
					fmt.Sprintf("0x%s", slot),
				); len(ids) > 0 {
					altMempoolIds = append(altMempoolIds, ids...)
				} else {
					return altMempoolIds, fmt.Errorf(
						"%s has forbidden %s to %s slot %s",
						v.EntityName,
						mode,
						addr2KnownEntity(v.Op, addr),
						slot,
					)
				}
			}
		}

		if mustStakeSlot != "" && !v.EntityIsStaked {
			return altMempoolIds, fmt.Errorf(
				"unstaked %s accessed %s slot %s",
				v.EntityName,
				addr2KnownEntity(v.Op, addr),
				mustStakeSlot,
			)
		}
	}

	return altMempoolIds, nil
}
