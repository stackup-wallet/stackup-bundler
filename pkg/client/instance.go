package client

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type Instance struct {
	ethClient            *ethclient.Client
	supportedEntryPoints []string
}

// Implements the method call for eth_sendUserOperation.
// Returns true if userOp was accepted otherwise returns an error.
func (i *Instance) Eth_sendUserOperation(op map[string]interface{}, ep string) (bool, error) {
	for _, v := range i.supportedEntryPoints {
		if common.HexToAddress(v) != common.HexToAddress(ep) {
			return false, errors.New("entryPoint: Implementation not supported")
		}
	}
	epAddr := common.HexToAddress(ep)
	entryPoint, err := entrypoint.NewEntrypoint(epAddr, i.ethClient)
	if err != nil {
		return false, err
	}

	userop, err := userop.New(op)
	if err != nil {
		return false, err
	}
	if err := userop.CheckSender(i.ethClient); err != nil {
		return false, err
	}
	if err := userop.CheckVerificationGasLimits(i.ethClient); err != nil {
		return false, err
	}
	if err := userop.CheckPaymasterAndData(i.ethClient, entryPoint); err != nil {
		return false, err
	}
	if err := userop.CheckCallGasLimit(i.ethClient); err != nil {
		return false, err
	}
	if err := userop.CheckFeePerGas(i.ethClient); err != nil {
		return false, err
	}
	if err := userop.CheckDuplicate(i.ethClient); err != nil {
		return false, err
	}

	return true, nil
}

// Implements the method call for eth_supportedEntryPoints.
// It returns the array of EntryPoint addresses that is supported by the client.
func (i *Instance) Eth_supportedEntryPoints() ([]string, error) {
	return i.supportedEntryPoints, nil
}

// Initializes a new ERC-4337 client with an ethClient instance
// and an array of supported EntryPoint addresses.
// The first address in the array is the preferred EntryPoint.
func New(ethClient *ethclient.Client, supportedEntryPoints []string) Instance {
	return Instance{
		ethClient:            ethClient,
		supportedEntryPoints: supportedEntryPoints,
	}
}
