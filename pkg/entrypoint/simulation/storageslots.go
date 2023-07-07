package simulation

import (
	"fmt"
	"math/big"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type storageSlots mapset.Set[string]

type storageSlotsByEntity map[common.Address]storageSlots

func newStorageSlotsByEntity(stakes EntityStakes, keccak []string) storageSlotsByEntity {
	storageSlotsByEntity := make(storageSlotsByEntity)

	for _, k := range keccak {
		value := common.Bytes2Hex(crypto.Keccak256(common.Hex2Bytes(k[2:])))

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
	Op         *userop.UserOperation
	EntryPoint common.Address

	// Parameters of specific entities required for all validation
	SenderSlots     storageSlots
	FactoryIsStaked bool

	// Parameters of the entity under validation
	EntityName     string
	EntityAddr     common.Address
	EntityAccess   tracer.AccessMap
	EntitySlots    storageSlots
	EntityIsStaked bool
}

func isAssociatedWith(slots storageSlots, slot string) bool {
	slotN, _ := big.NewInt(0).SetString(fmt.Sprintf("0x%s", slot), 0)
	for _, k := range slots.ToSlice() {
		kn, _ := big.NewInt(0).SetString(fmt.Sprintf("0x%s", k), 0)
		if slotN.Cmp(kn) >= 0 && slotN.Cmp(big.NewInt(0).Add(kn, big.NewInt(128))) <= 0 {
			return true
		}
	}

	return false
}

func (v *storageSlotsValidator) Process() error {
	senderSlots := v.SenderSlots
	if senderSlots == nil {
		senderSlots = mapset.NewSet[string]()
	}
	entitySlots := v.EntitySlots
	if entitySlots == nil {
		entitySlots = mapset.NewSet[string]()
	}

	for addr, access := range v.EntityAccess {
		if addr == v.Op.Sender || addr == v.EntryPoint {
			continue
		}

		var mustStakeSlot string
		accessTypes := map[string]tracer.Counts{
			"read":  access.Reads,
			"write": access.Writes,
		}
		for key, slotCount := range accessTypes {
			for slot := range slotCount {
				if isAssociatedWith(senderSlots, slot) {
					if len(v.Op.InitCode) > 0 && !v.FactoryIsStaked {
						mustStakeSlot = slot
					} else {
						continue
					}
				} else if isAssociatedWith(entitySlots, slot) || addr == v.EntityAddr {
					mustStakeSlot = slot
				} else {
					return fmt.Errorf("%s has forbidden %s to %s slot %s", v.EntityName, key, addr2KnownEntity(v.Op, addr), slot)
				}
			}
		}

		if mustStakeSlot != "" && !v.EntityIsStaked {
			return fmt.Errorf(
				"unstaked %s accessed %s slot %s",
				v.EntityName,
				addr2KnownEntity(v.Op, addr),
				mustStakeSlot,
			)
		}
	}

	return nil
}
