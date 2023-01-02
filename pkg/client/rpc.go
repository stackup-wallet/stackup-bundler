package client

import (
	"errors"

	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
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

// Eth_sendUserOperation routes method calls to *Client.SendUserOperation.
func (r *RpcAdapter) Eth_sendUserOperation(op map[string]any, ep string) (string, error) {
	return r.client.SendUserOperation(op, ep)
}

// Eth_estimateUserOperationGas routes method calls to *Client.EstimateUserOperationGas.
func (r *RpcAdapter) Eth_estimateUserOperationGas(op map[string]any, ep string) (*gas.GasEstimates, error) {
	return r.client.EstimateUserOperationGas(op, ep)
}

// Eth_getUserOperationReceipt routes method calls to *Client.GetUserOperationReceipt.
func (r *RpcAdapter) Eth_getUserOperationReceipt(
	userOpHash string,
) (*entrypoint.UserOperationReceipt, error) {
	return r.client.GetUserOperationReceipt(userOpHash)
}

// Eth_supportedEntryPoints routes method calls to *Client.SupportedEntryPoints.
func (r *RpcAdapter) Eth_supportedEntryPoints() ([]string, error) {
	return r.client.SupportedEntryPoints()
}

// Eth_chainId routes method calls to *Client.ChainID.
func (r *RpcAdapter) Eth_chainId() (string, error) {
	return r.client.ChainID()
}

// Debug_bundler_clearState routes method calls to *Debug.ClearState.
func (r *RpcAdapter) Debug_bundler_clearState() (string, error) {
	if r.debug == nil {
		return "", errors.New("rpc: debug mode is not enabled")
	}

	return r.debug.ClearState()
}

// Debug_bundler_dumpMempool routes method calls to *Debug.DumpMempool.
func (r *RpcAdapter) Debug_bundler_dumpMempool(ep string) ([]map[string]any, error) {
	if r.debug == nil {
		return []map[string]any{}, errors.New("rpc: debug mode is not enabled")
	}

	return r.debug.DumpMempool(ep)
}

// Debug_bundler_sendBundleNow routes method calls to *Debug.SendBundleNow.
func (r *RpcAdapter) Debug_bundler_sendBundleNow() (string, error) {
	if r.debug == nil {
		return "", errors.New("rpc: debug mode is not enabled")
	}

	return r.debug.SendBundleNow()
}

// Debug_bundler_setBundlingMode routes method calls to *Debug.SetBundlingMode.
func (r *RpcAdapter) Debug_bundler_setBundlingMode(mode string) (string, error) {
	if r.debug == nil {
		return "", errors.New("rpc: debug mode is not enabled")
	}

	return r.debug.SetBundlingMode(mode)
}
