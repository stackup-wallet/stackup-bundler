package client

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Client controls the end to end process of adding incoming UserOperations to the mempool.
type Client struct {
	mempool              *mempool.Interface
	chainID              *big.Int
	supportedEntryPoints []common.Address
	userOpHandler        modules.UserOpHandlerFunc
}

func (i *Client) parseEntryPointAddress(ep string) (common.Address, error) {
	for _, addr := range i.supportedEntryPoints {
		if common.HexToAddress(ep) == addr {
			return addr, nil
		}
	}

	return common.Address{}, errors.New("entryPoint: Implementation not supported")
}

// UseModules defines the UserOpHandlers to process a userOp after it has gone through the standard checks.
func (i *Client) UseModules(handlers ...modules.UserOpHandlerFunc) {
	i.userOpHandler = modules.ComposeUserOpHandlerFunc(handlers...)
}

// SendUserOperation implements the method call for eth_sendUserOperation.
// It returns true if userOp was accepted otherwise returns an error.
func (i *Client) SendUserOperation(op map[string]any, ep string) (bool, error) {
	// Check EntryPoint and userOp is valid.
	epAddr, err := i.parseEntryPointAddress(ep)
	if err != nil {
		return false, err
	}
	userOp, err := userop.New(op)
	if err != nil {
		return false, err
	}

	// Check mempool for duplicates and only replace under the following circumstances:
	//
	//	1. the nonce remains the same
	//	2. the new maxPriorityFeePerGas is higher
	//	3. the new maxFeePerGas is increased equally
	memOp, err := i.mempool.GetOp(epAddr, userOp.Sender)
	if err != nil {
		return false, err
	}
	if memOp != nil {
		if memOp.Nonce.Cmp(memOp.Nonce) != 0 {
			return false, errors.New("sender: Has userOp in mempool with a different nonce")
		}

		if memOp.MaxPriorityFeePerGas.Cmp(memOp.MaxPriorityFeePerGas) <= 0 {
			return false, errors.New("sender: Has userOp in mempool with same or higher priority fee")
		}

		diff := big.NewInt(0)
		mf := big.NewInt(0)
		diff.Sub(memOp.MaxPriorityFeePerGas, memOp.MaxPriorityFeePerGas)
		if memOp.MaxFeePerGas.Cmp(mf.Add(memOp.MaxFeePerGas, diff)) != 0 {
			return false, errors.New("sender: Replaced userOp must have an equally higher max fee")
		}
	}

	// Run through client module stack.
	ctx := modules.NewUserOpHandlerContext(userOp, epAddr, i.chainID)
	if err := i.userOpHandler(ctx); err != nil {
		return false, err
	}

	// Add userOp to mempool.
	if err := i.mempool.AddOp(epAddr, ctx.UserOp); err != nil {
		return false, err
	}
	return true, nil
}

// SupportedEntryPoints implements the method call for eth_supportedEntryPoints.
// It returns the array of EntryPoint addresses that is supported by the client.
func (i *Client) SupportedEntryPoints() ([]string, error) {
	slc := []string{}
	for _, ep := range i.supportedEntryPoints {
		slc = append(slc, ep.String())
	}

	return slc, nil
}
