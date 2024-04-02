package entrypoint

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	entrypointV06 "github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/bindings/v06"
)
import entrypointV07 "github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/bindings/v07"

type StakeManagerDepositInfoFactoryFunc func() IStakeManagerDepositInfo

var StakeManagerDepositInfoFactories = make(map[string]StakeManagerDepositInfoFactoryFunc)

type IStakeManagerDepositInfo interface {
	GetVersion() string
}

func NewStakeManagerByVersion(version string) IStakeManagerDepositInfo {
	if factory, ok := StakeManagerDepositInfoFactories[version]; ok {
		return factory()
	}
	return nil
}

// NewEntrypoint creates a new instance of Entrypoint, bound to a specific deployed contract.
func NewEntrypoint(address common.Address, backend bind.ContractBackend) (*entrypointV06.Entrypoint, error) {
	contract, err := bindEntrypoint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Entrypoint{EntrypointCaller: EntrypointCaller{contract: contract}, EntrypointTransactor: EntrypointTransactor{contract: contract}, EntrypointFilterer: EntrypointFilterer{contract: contract}}, nil
}

func init() {
	StakeManagerDepositInfoFactories[entrypointV06.VERSION] = func() IStakeManagerDepositInfo {
		return &entrypointV06.IStakeManagerDepositInfo{}
	}
	StakeManagerDepositInfoFactories[entrypointV07.VERSION] = func() IStakeManagerDepositInfo {
		return &entrypointV07.IStakeManagerDepositInfo{}
	}
}
