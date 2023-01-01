package client

import (
	"errors"
)

// RpcAdapter is an adapter for routing JSON-RPC method calls to the correct client functions.
type RpcAdapter struct {
	client *Client
	debug  *Debug
}

// NewRpcAdapter initializes a new RpcAdapter which can be used with a JSON-RPC server.
func NewRpcAdapter(client *Client, debug *Debug) *RpcAdapter {
	return &RpcAdapter{client, debug}
}

// Eth_sendUserOperation routes eth_sendUserOperation method calls to *Client.SendUserOperation.
func (r *RpcAdapter) Eth_sendUserOperation(op map[string]any, ep string) (string, error) {
	return r.client.SendUserOperation(op, ep)
}

// Eth_supportedEntryPoints routes eth_supportedEntryPoints method calls to *Client.SupportedEntryPoints.
func (r *RpcAdapter) Eth_supportedEntryPoints() ([]string, error) {
	return r.client.SupportedEntryPoints()
}

// Eth_chainId routes eth_chainId method calls to *Client.ChainID.
func (r *RpcAdapter) Eth_chainId() (string, error) {
	return r.client.ChainID()
}

// Debug_bundler_sendBundleNow routes eth_chainId method calls to *Debug.SendBundleNow.
func (r *RpcAdapter) Debug_bundler_sendBundleNow() (string, error) {
	if r.debug == nil {
		return "", errors.New("rpc: debug mode is not enabled")
	}

	return r.debug.SendBundleNow()
}
