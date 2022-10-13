package client

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type Instance struct {
	ethClient            *ethclient.Client
	supportedEntryPoints []string
}

// Implements the method call for eth_sendUserOperation.
// Returns true if userOp was accepted otherwise returns an error.
func (i *Instance) Eth_sendUserOperation(op map[string]interface{}) (bool, error) {
	userop, err := userop.FromMap(op)
	if err != nil {
		return false, err
	}

	if err := userop.CheckSender(i.ethClient); err != nil {
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
