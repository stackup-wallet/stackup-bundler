package client

import (
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type Instance struct {
	supportedEntryPoints []string
}

// Implements the method call for eth_sendUserOperation.
// Returns true if userOp was accepted otherwise returns an error.
func (i *Instance) Eth_sendUserOperation(op map[string]interface{}) (bool, error) {
	_, err := userop.FromMap(op)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Implements the method call for eth_supportedEntryPoints.
// It returns the array of EntryPoint addresses that is supported by the client.
func (i *Instance) Eth_supportedEntryPoints() ([]string, error) {
	return i.supportedEntryPoints, nil
}

// Initializes a new ERC-4337 client with an array of supported EntryPoint addresses.
// The first address in the array is the preferred EntryPoint.
func New(supportedEntryPoints []string) Instance {
	return Instance{
		supportedEntryPoints: supportedEntryPoints,
	}
}
