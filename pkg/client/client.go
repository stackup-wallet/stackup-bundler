// Package client provides the mediator for processing incoming UserOperations to the bundler.
package client

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-logr/logr"
	"github.com/stackup-wallet/stackup-bundler/internal/logger"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/filter"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/noop"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Client controls the end to end process of adding incoming UserOperations to the mempool. It also
// implements the required RPC methods as specified in EIP-4337.
type Client struct {
	mempool               *mempool.Mempool
	chainID               *big.Int
	supportedEntryPoints  []common.Address
	userOpHandler         modules.UserOpHandlerFunc
	logger                logr.Logger
	getUserOpReceipt      GetUserOpReceiptFunc
	getSimulateValidation GetSimulateValidationFunc
	getCallGasEstimate    GetCallGasEstimateFunc
	getUserOpByHash       GetUserOpByHashFunc
}

// New initializes a new ERC-4337 client which can be extended with modules for validating UserOperations
// that are allowed to be added to the mempool.
func New(
	mempool *mempool.Mempool,
	chainID *big.Int,
	supportedEntryPoints []common.Address,
) *Client {
	return &Client{
		mempool:               mempool,
		chainID:               chainID,
		supportedEntryPoints:  supportedEntryPoints,
		userOpHandler:         noop.UserOpHandler,
		logger:                logger.NewZeroLogr().WithName("client"),
		getUserOpReceipt:      getUserOpReceiptNoop(),
		getSimulateValidation: getSimulateValidationNoop(),
		getCallGasEstimate:    getCallGasEstimateNoop(),
		getUserOpByHash:       getUserOpByHashNoop(),
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

// SetGetUserOpReceiptFunc defines a general function for fetching a UserOpReceipt given a userOpHash and
// EntryPoint address. This function is called in *Client.GetUserOperationReceipt.
func (i *Client) SetGetUserOpReceiptFunc(fn GetUserOpReceiptFunc) {
	i.getUserOpReceipt = fn
}

// SetGetSimulateValidationFunc defines a general function for fetching simulateValidation results given a
// userOp and EntryPoint address. This function is called in *Client.EstimateUserOperationGas.
func (i *Client) SetGetSimulateValidationFunc(fn GetSimulateValidationFunc) {
	i.getSimulateValidation = fn
}

// SetGetCallGasEstimateFunc defines a general function for fetching an estimate for callGasLimit given a
// userOp and EntryPoint address. This function is called in *Client.EstimateUserOperationGas.
func (i *Client) SetGetCallGasEstimateFunc(fn GetCallGasEstimateFunc) {
	i.getCallGasEstimate = fn
}

// SetGetUserOpByHashFunc defines a general function for fetching a userOp given a userOpHash, EntryPoint
// address, and chain ID. This function is called in *Client.GetUserOperationByHash.
func (i *Client) SetGetUserOpByHashFunc(fn GetUserOpByHashFunc) {
	i.getUserOpByHash = fn
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

// EstimateUserOperationGas returns estimates for PreVerificationGas, VerificationGas, and CallGasLimit given
// a UserOperation and EntryPoint address. The signature field and current gas values will not be validated
// although there should be dummy values in place for the most reliable results (e.g. a signature with the
// correct length).
func (i *Client) EstimateUserOperationGas(op map[string]any, ep string) (*gas.GasEstimates, error) {
	// Init logger
	l := i.logger.WithName("eth_estimateUserOperationGas")

	// Check EntryPoint and userOp is valid.
	epAddr, err := i.parseEntryPointAddress(ep)
	if err != nil {
		l.Error(err, "eth_estimateUserOperationGas error")
		return nil, err
	}
	l = l.
		WithValues("entrypoint", epAddr.String()).
		WithValues("chain_id", i.chainID.String())

	userOp, err := userop.New(op)
	if err != nil {
		l.Error(err, "eth_estimateUserOperationGas error")
		return nil, err
	}
	hash := userOp.GetUserOpHash(epAddr, i.chainID)
	l = l.WithValues("userop_hash", hash)

	sim, err := i.getSimulateValidation(epAddr, userOp)
	if err != nil {
		l.Error(err, "eth_estimateUserOperationGas error")
		return nil, err
	}

	cg, err := i.getCallGasEstimate(epAddr, userOp)
	if err != nil {
		l.Error(err, "eth_estimateUserOperationGas error")
		return nil, err
	}

	l.Info("eth_estimateUserOperationGas ok")
	return &gas.GasEstimates{
		PreVerificationGas: gas.NewDefaultOverhead().CalcPreVerificationGas(userOp),
		VerificationGas:    sim.ReturnInfo.PreOpGas,
		CallGasLimit:       big.NewInt(int64(cg)),
	}, nil
}

// GetUserOperationReceipt fetches a UserOperation receipt based on a userOpHash returned by
// *Client.SendUserOperation.
func (i *Client) GetUserOperationReceipt(
	hash string,
) (*filter.UserOperationReceipt, error) {
	// Init logger
	l := i.logger.WithName("eth_getUserOperationReceipt").WithValues("userop_hash", hash)

	ev, err := i.getUserOpReceipt(hash, i.supportedEntryPoints[0])
	if err != nil {
		l.Error(err, "eth_getUserOperationReceipt error")
		return nil, err
	}

	l.Info("eth_getUserOperationReceipt ok")
	return ev, nil
}

// GetUserOperationByHash returns a UserOperation based on a given userOpHash returned by
// *Client.SendUserOperation.
func (i *Client) GetUserOperationByHash(hash string) (*filter.HashLookupResult, error) {
	// Init logger
	l := i.logger.WithName("eth_getUserOperationByHash").WithValues("userop_hash", hash)

	res, err := i.getUserOpByHash(hash, i.supportedEntryPoints[0], i.chainID)
	if err != nil {
		l.Error(err, "eth_getUserOperationByHash error")
		return nil, err
	}

	return res, nil
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
	return hexutil.EncodeBig(i.chainID), nil
}
