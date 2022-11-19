package client

// RpcAdapter is an adapter for routing JSON-RPC method calls to the correct client functions.
type RpcAdapter struct {
	client *Client
}

// NewRpcAdapter initializes a new RpcAdapter which can be used with a JSON-RPC server.
func NewRpcAdapter(client *Client) *RpcAdapter {
	return &RpcAdapter{client}
}

// Eth_sendUserOperation routes eth_sendUserOperation method calls to client.SendUserOperation.
func (r *RpcAdapter) Eth_sendUserOperation(op map[string]any, ep string) (bool, error) {
	return r.client.SendUserOperation(op, ep)
}

// Eth_supportedEntryPoints routes eth_supportedEntryPoints method calls to client.supportedEntryPoints.
func (r *RpcAdapter) Eth_supportedEntryPoints() ([]string, error) {
	return r.client.SupportedEntryPoints()
}
