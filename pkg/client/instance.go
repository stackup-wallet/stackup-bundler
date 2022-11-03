package client

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type Instance struct {
	ethClient            *ethclient.Client
	mempool              *mempool.Interface
	chainID              *big.Int
	supportedEntryPoints []common.Address
}

func (i *Instance) parseEntryPointAddress(ep string) (common.Address, error) {
	for _, addr := range i.supportedEntryPoints {
		if common.HexToAddress(ep) == addr {
			return addr, nil
		}
	}

	return common.Address{}, errors.New("entryPoint: Implementation not supported")
}

// Implements the method call for eth_sendUserOperation.
// Returns true if userOp was accepted otherwise returns an error.
func (i *Instance) Eth_sendUserOperation(op map[string]interface{}, ep string) (bool, error) {
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

	// Add to mempool
	err = i.mempool.AddOp(epAddr, userop)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Implements the method call for eth_supportedEntryPoints.
// It returns the array of EntryPoint addresses that is supported by the client.
func (i *Instance) Eth_supportedEntryPoints() ([]string, error) {
	slc := []string{}
	for _, ep := range i.supportedEntryPoints {
		slc = append(slc, ep.String())
	}

	return slc, nil
}
