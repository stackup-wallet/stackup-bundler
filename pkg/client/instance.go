package client

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type Instance struct {
	ethClient            *ethclient.Client
	mempool              *mempool.ClientInterface
	chainID              *big.Int
	supportedEntryPoints []string
}

func (i *Instance) parseEntryPointAddress(ep string) (common.Address, error) {
	for _, v := range i.supportedEntryPoints {
		addr := common.HexToAddress(v)
		if addr == common.HexToAddress(ep) {
			return addr, nil
		}
	}

	return common.Address{}, errors.New("entryPoint: Implementation not supported")
}

// Implements the method call for eth_sendUserOperation.
// Returns true if userOp was accepted otherwise returns an error.
func (i *Instance) Eth_sendUserOperation(op map[string]interface{}, ep string) (bool, error) {
	epAddr, err := i.parseEntryPointAddress(ep)
	if err != nil {
		return false, err
	}
	entryPoint, err := entrypoint.NewEntrypoint(epAddr, i.ethClient)
	if err != nil {
		return false, err
	}

	// Run sanity checks
	userop, err := userop.New(op)
	if err != nil {
		return false, err
	}
	if err := checkSender(userop, i.ethClient); err != nil {
		return false, err
	}
	if err := checkVerificationGasLimits(userop, i.ethClient); err != nil {
		return false, err
	}
	if err := checkPaymasterAndData(userop, i.ethClient, entryPoint); err != nil {
		return false, err
	}
	if err := checkCallGasLimit(userop, i.ethClient); err != nil {
		return false, err
	}
	if err := checkFeePerGas(userop, i.ethClient); err != nil {
		return false, err
	}
	if err := checkDuplicates(userop, i.mempool); err != nil {
		return false, err
	}

	// Add to mempool
	res, err := i.mempool.Add(userop.Sender.String(), userop, epAddr)
	if err != nil {
		return false, err
	}

	return res, nil
}

// Implements the method call for eth_supportedEntryPoints.
// It returns the array of EntryPoint addresses that is supported by the client.
func (i *Instance) Eth_supportedEntryPoints() ([]string, error) {
	return i.supportedEntryPoints, nil
}
