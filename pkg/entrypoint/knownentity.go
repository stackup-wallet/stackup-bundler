package entrypoint

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type knownEntity map[string]struct {
	Address  common.Address
	Info     tracer.NumberLevelInfo
	IsStaked bool
}

func newKnownEntity(
	op *userop.UserOperation,
	res *tracer.BundlerCollectorReturn,
	stakes EntityStakesMap,
) (knownEntity, error) {
	if len(res.NumberLevels) != 3 {
		return nil, fmt.Errorf("unexpected NumberLevels length in tracing result: %d", len(res.NumberLevels))
	}

	return knownEntity{
		"account": {
			Address:  op.Sender,
			Info:     res.NumberLevels[accountNumberLevel],
			IsStaked: stakes[op.Sender] != nil && stakes[op.Sender].Staked,
		},
		"factory": {
			Address:  op.GetFactory(),
			Info:     res.NumberLevels[factoryNumberLevel],
			IsStaked: stakes[op.GetFactory()] != nil && stakes[op.GetFactory()].Staked,
		},
		"paymaster": {
			Address:  op.GetPaymaster(),
			Info:     res.NumberLevels[paymasterNumberLevel],
			IsStaked: stakes[op.GetPaymaster()] != nil && stakes[op.GetPaymaster()].Staked,
		},
	}, nil
}

func addr2KnownEntity(op *userop.UserOperation, addr common.Address) string {
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
