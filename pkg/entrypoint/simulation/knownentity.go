package simulation

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/methods"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type knownEntity map[string]struct {
	Address  common.Address
	Info     tracer.CallFromEntryPointInfo
	IsStaked bool
}

func newKnownEntity(
	op *userop.UserOperation,
	entryPoint common.Address,
	res *tracer.BundlerCollectorReturn,
	stakes EntityStakes,
) (knownEntity, error) {
	si := tracer.CallFromEntryPointInfo{}
	fi := tracer.CallFromEntryPointInfo{}
	pi := tracer.CallFromEntryPointInfo{}
	for _, c := range res.CallsFromEntryPoint {
		switch c.TopLevelTargetAddress {
		case op.Sender:
			si = c
		case op.GetPaymaster():
			pi = c
		default:
			if c.TopLevelMethodSig.String() == methods.CreateSenderSelector {
				fi = c
			}
		}
	}

	return knownEntity{
		"account": {
			Address:  op.Sender,
			Info:     si,
			IsStaked: stakes[op.Sender] != nil && stakes[op.Sender].Staked,
		},
		"factory": {
			Address:  op.GetFactory(),
			Info:     fi,
			IsStaked: stakes[op.GetFactory()] != nil && stakes[op.GetFactory()].Staked,
		},
		"paymaster": {
			Address:  op.GetPaymaster(),
			Info:     pi,
			IsStaked: stakes[op.GetPaymaster()] != nil && stakes[op.GetPaymaster()].Staked,
		},
	}, nil
}

func addr2KnownEntity(op *userop.UserOperation, addr common.Address) string {
	if addr == op.GetFactory() {
		return "factory"
	} else if addr == op.Sender {
		return "account"
	} else if addr == op.GetPaymaster() {
		return "paymaster"
	} else {
		return addr.String()
	}
}
