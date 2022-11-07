package client

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Instance is a representation of an ERC-4337 client.
type Instance struct {
	ethClient            *ethclient.Client
	mempool              *mempool.Interface
	chainID              *big.Int
	supportedEntryPoints []common.Address
	userOpHandler        modules.UserOpHandlerFunc
}

func (i *Instance) parseEntryPointAddress(ep string) (common.Address, error) {
	for _, addr := range i.supportedEntryPoints {
		if common.HexToAddress(ep) == addr {
			return addr, nil
		}
	}

	return common.Address{}, errors.New("entryPoint: Implementation not supported")
}

// UseModules defines the UserOpHandlers to process a userOp after it has gone through the standard checks.
func (i *Instance) UseModules(handlers ...modules.UserOpHandlerFunc) {
	i.userOpHandler = modules.ComposeUserOpHandlerFunc(handlers...)
}

// SendUserOperation implements the method call for eth_sendUserOperation.
// It returns true if userOp was accepted otherwise returns an error.
func (i *Instance) SendUserOperation(op map[string]any, ep string) (bool, error) {
	epAddr, err := i.parseEntryPointAddress(ep)
	if err != nil {
		return false, err
	}
	entryPoint, err := entrypoint.NewEntrypoint(epAddr, i.ethClient)
	if err != nil {
		return false, err
	}

	// Run sanity checks
	userop, err := userop.New(op)
	if err != nil {
		return false, err
	}
	if err := checkSender(userop, i.ethClient); err != nil {
		return false, err
	}
	if err := checkVerificationGasLimits(userop, i.ethClient); err != nil {
		return false, err
	}
	if err := checkPaymasterAndData(userop, i.ethClient, entryPoint); err != nil {
		return false, err
	}
	if err := checkCallGasLimit(userop, i.ethClient); err != nil {
		return false, err
	}
	if err := checkFeePerGas(userop, i.ethClient); err != nil {
		return false, err
	}
	if err := checkDuplicates(userop, epAddr, i.mempool); err != nil {
		return false, err
	}

	// Run through additional modules
	ctx := modules.NewUserOpHandlerContext(userop)
	if err := i.userOpHandler(ctx); err != nil {
		return false, err
	}

	// Add to mempool
	err = i.mempool.AddOp(epAddr, ctx.UserOp)
	if err != nil {
		return false, err
	}

	return true, nil
}

// SupportedEntryPoints implements the method call for eth_supportedEntryPoints.
// It returns the array of EntryPoint addresses that is supported by the client.
func (i *Instance) SupportedEntryPoints() ([]string, error) {
	slc := []string{}
	for _, ep := range i.supportedEntryPoints {
		slc = append(slc, ep.String())
	}

	return slc, nil
}
