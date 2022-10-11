package client

type Instance struct {
	supportedEntryPoints []string
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
