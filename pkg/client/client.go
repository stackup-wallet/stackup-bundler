// Package client provides the mediator for processing incoming UserOperations to the bundler.
package client

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-logr/logr"
	"github.com/stackup-wallet/stackup-bundler/internal/logger"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/noop"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Client controls the end to end process of adding incoming UserOperations to the mempool. It also
// implements the required RPC methods as specified in EIP-4337.
type Client struct {
	mempool              *mempool.Mempool
	chainID              *big.Int
	supportedEntryPoints []common.Address
	userOpHandler        modules.UserOpHandlerFunc
	logger               logr.Logger
}

// New initializes a new ERC-4337 client which can be extended with modules for validating UserOperations
// that are allowed to be added to the mempool.
func New(mempool *mempool.Mempool, chainID *big.Int, supportedEntryPoints []common.Address) *Client {
	return &Client{
		mempool:              mempool,
		chainID:              chainID,
		supportedEntryPoints: supportedEntryPoints,
		userOpHandler:        noop.UserOpHandler,
		logger:               logger.NewZeroLogr().WithName("client"),
	}
}

func (i *Client) parseEntryPointAddress(ep string) (common.Address, error) {
	for _, addr := range i.supportedEntryPoints {
		if common.HexToAddress(ep) == addr {
			return addr, nil
		}
	}

	return common.Address{}, errors.New("entryPoint: Implementation not supported")
}

// UseLogger defines the logger object used by the Client instance based on the go-logr/logr interface.
func (i *Client) UseLogger(logger logr.Logger) {
	i.logger = logger.WithName("client")
}

// UseModules defines the UserOpHandlers to process a userOp after it has gone through the standard checks.
func (i *Client) UseModules(handlers ...modules.UserOpHandlerFunc) {
	i.userOpHandler = modules.ComposeUserOpHandlerFunc(handlers...)
}

// SendUserOperation implements the method call for eth_sendUserOperation.
// It returns true if userOp was accepted otherwise returns an error.
func (i *Client) SendUserOperation(op map[string]any, ep string) (string, error) {
	// Init logger
	l := i.logger.WithName("eth_sendUserOperation")

	// Check EntryPoint and userOp is valid.
	epAddr, err := i.parseEntryPointAddress(ep)
	if err != nil {
		l.Error(err, "eth_sendUserOperation error")
		return "", err
	}
	l = l.
		WithValues("entrypoint", epAddr.String()).
		WithValues("chain_id", i.chainID.String())

	userOp, err := userop.New(op)
	if err != nil {
		l.Error(err, "eth_sendUserOperation error")
		return "", err
	}
	hash := userOp.GetUserOpHash(epAddr, i.chainID)
	l = l.WithValues("userop_hash", hash)

	// Fetch any pending UserOperations in the mempool by the same sender
	penOps, err := i.mempool.GetOps(epAddr, userOp.Sender)
	if err != nil {
		l.Error(err, "eth_sendUserOperation error")
		return "", err
	}

	// Run through client module stack.
	ctx := modules.NewUserOpHandlerContext(userOp, penOps, epAddr, i.chainID)
	if err := i.userOpHandler(ctx); err != nil {
		l.Error(err, "eth_sendUserOperation error")
		return "", err
	}

	// Add userOp to mempool.
	if err := i.mempool.AddOp(epAddr, ctx.UserOp); err != nil {
		l.Error(err, "eth_sendUserOperation error")
		return "", err
	}

	l.Info("eth_sendUserOperation ok")
	return hash.String(), nil
}

// SupportedEntryPoints implements the method call for eth_supportedEntryPoints. It returns the array of
// EntryPoint addresses that is supported by the client. The first address in the array is the preferred
// EntryPoint.
func (i *Client) SupportedEntryPoints() ([]string, error) {
	slc := []string{}
	for _, ep := range i.supportedEntryPoints {
		slc = append(slc, ep.String())
	}

	return slc, nil
}

// ChainID implements the method call for eth_chainId. It returns the current chainID used by the client.
// This method is used to validate that the client's chainID is in sync with the caller.
func (i *Client) ChainID() (string, error) {
	return i.chainID.String(), nil
}
