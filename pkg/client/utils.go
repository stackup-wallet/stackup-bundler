package client

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/filter"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/reverts"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/simulation"
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

// GetSimulateValidationFunc is a general interface for fetching simulateValidation results given a userOp
// and EntryPoint address.
type GetSimulateValidationFunc = func(ep common.Address, op *userop.UserOperation) (*reverts.ValidationResultRevert, error)

func getSimulateValidationNoop() GetSimulateValidationFunc {
	return func(ep common.Address, op *userop.UserOperation) (*reverts.ValidationResultRevert, error) {
		//lint:ignore ST1005 This needs to match the bundler test spec.
		return nil, errors.New("Missing/invalid userOpHash")
	}
}

// GetSimulateValidationWithRpcClient returns an implementation of GetSimulateValidationFunc that relies on a
// rpc client to fetch simulateValidation results.
func GetSimulateValidationWithRpcClient(rpc *rpc.Client) GetSimulateValidationFunc {
	return func(ep common.Address, op *userop.UserOperation) (*reverts.ValidationResultRevert, error) {
		return simulation.SimulateValidation(rpc, ep, op)
	}
}

// GetCallGasEstimateFunc is a general interface for fetching an estimate for callGasLimit given a userOp and
// EntryPoint address.
type GetCallGasEstimateFunc = func(ep common.Address, op *userop.UserOperation) (uint64, error)

func getCallGasEstimateNoop() GetCallGasEstimateFunc {
	return func(ep common.Address, op *userop.UserOperation) (uint64, error) {
		//lint:ignore ST1005 This needs to match the bundler test spec.
		return 0, errors.New("Missing/invalid userOpHash")
	}
}

// GetCallGasEstimateWithEthClient returns an implementation of GetCallGasEstimateFunc that relies on an eth
// client to fetch an estimate for callGasLimit.
func GetCallGasEstimateWithEthClient(eth *ethclient.Client) GetCallGasEstimateFunc {
	return func(ep common.Address, op *userop.UserOperation) (uint64, error) {
		return gas.CallGasEstimate(eth, ep, op)
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
