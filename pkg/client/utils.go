package client

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/filter"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// GetUserOpReceiptFunc is a general interface for fetching a UserOperationReceipt given a userOpHash and
// EntryPoint address.
type GetUserOpReceiptFunc = func(hash string, ep common.Address) (*filter.UserOperationReceipt, error)

func getUserOpReceiptNoop() GetUserOpReceiptFunc {
	return func(hash string, ep common.Address) (*filter.UserOperationReceipt, error) {
		//lint:ignore ST1005 This needs to match the bundler test spec.
		return nil, errors.New("Missing/invalid userOpHash")
	}
}

// GetUserOpReceiptWithEthClient returns an implementation of GetUserOpReceiptFunc that relies on an eth
// client to fetch a UserOperationReceipt.
func GetUserOpReceiptWithEthClient(eth *ethclient.Client) GetUserOpReceiptFunc {
	return func(hash string, ep common.Address) (*filter.UserOperationReceipt, error) {
		return filter.GetUserOperationReceipt(eth, hash, ep)
	}
}

// GetGasEstimateFunc is a general interface for fetching an estimate for verificationGasLimit and
// callGasLimit given a userOp and EntryPoint address.
type GetGasEstimateFunc = func(ep common.Address, op *userop.UserOperation) (verificationGas uint64, callGas uint64, err error)

func getGasEstimateNoop() GetGasEstimateFunc {
	return func(ep common.Address, op *userop.UserOperation) (verificationGas uint64, callGas uint64, err error) {
		//lint:ignore ST1005 This needs to match the bundler test spec.
		return 0, 0, errors.New("Missing/invalid userOpHash")
	}
}

// GetGasEstimateWithEthClient returns an implementation of GetGasEstimateFunc that relies on an eth client to
// fetch an estimate for verificationGasLimit and callGasLimit.
func GetGasEstimateWithEthClient(
	rpc *rpc.Client,
	ov *gas.Overhead,
	chain *big.Int,
	maxGasLimit *big.Int,
) GetGasEstimateFunc {
	return func(ep common.Address, op *userop.UserOperation) (verificationGas uint64, callGas uint64, err error) {
		return gas.EstimateGas(&gas.EstimateInput{
			Rpc:         rpc,
			EntryPoint:  ep,
			Op:          op,
			Ov:          ov,
			ChainID:     chain,
			MaxGasLimit: maxGasLimit,
		})
	}
}

// GetUserOpByHashFunc is a general interface for fetching a UserOperation given a userOpHash, EntryPoint
// address, and chain ID.
type GetUserOpByHashFunc func(hash string, ep common.Address, chain *big.Int) (*filter.HashLookupResult, error)

func getUserOpByHashNoop() GetUserOpByHashFunc {
	return func(hash string, ep common.Address, chain *big.Int) (*filter.HashLookupResult, error) {
		//lint:ignore ST1005 This needs to match the bundler test spec.
		return nil, errors.New("Missing/invalid userOpHash")
	}
}

// GetUserOpByHashWithEthClient returns an implementation of GetUserOpByHashFunc that relies on an eth client
// to fetch a UserOperation.
func GetUserOpByHashWithEthClient(eth *ethclient.Client) GetUserOpByHashFunc {
	return func(hash string, ep common.Address, chain *big.Int) (*filter.HashLookupResult, error) {
		return filter.GetUserOperationByHash(eth, hash, ep, chain)
	}
}
