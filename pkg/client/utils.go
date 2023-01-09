package client

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// GetUserOpReceiptFunc is a general interface for fetching a UserOperationReceipt given a userOpHash and
// EntryPoint address.
type GetUserOpReceiptFunc = func(hash string, ep common.Address) (*entrypoint.UserOperationReceipt, error)

func getUserOpReceiptNoop() GetUserOpReceiptFunc {
	return func(hash string, ep common.Address) (*entrypoint.UserOperationReceipt, error) {
		//lint:ignore ST1005 This needs to match the bundler test spec.
		return nil, errors.New("Missing/invalid userOpHash")
	}
}

// GetUserOpReceiptWithEthClient returns an implementation of GetUserOpReceiptFunc that relies on an eth
// client to fetch a UserOperationReceipt.
func GetUserOpReceiptWithEthClient(eth *ethclient.Client) GetUserOpReceiptFunc {
	return func(hash string, ep common.Address) (*entrypoint.UserOperationReceipt, error) {
		return entrypoint.GetUserOperationReceipt(eth, hash, ep)
	}
}

// GetSimulateValidationFunc is a general interface for fetching simulate validation results given a userOp
// and EntryPoint address.
type GetSimulateValidationFunc = func(ep common.Address, op *userop.UserOperation) (*entrypoint.ValidationResultRevert, error)

func getSimulateValidationNoop() GetSimulateValidationFunc {
	return func(ep common.Address, op *userop.UserOperation) (*entrypoint.ValidationResultRevert, error) {
		//lint:ignore ST1005 This needs to match the bundler test spec.
		return nil, errors.New("Missing/invalid userOpHash")
	}
}

// GetSimulateValidationWithRpcClient returns an implementation of GetSimulateValidationFunc that relies on an
// rpc client to fetch simulate validation results.
func GetSimulateValidationWithRpcClient(rpc *rpc.Client) GetSimulateValidationFunc {
	return func(ep common.Address, op *userop.UserOperation) (*entrypoint.ValidationResultRevert, error) {
		return entrypoint.SimulateValidation(rpc, ep, op)
	}
}

type GetCallGasEstimateFunc = func(ep common.Address, op *userop.UserOperation) (uint64, error)

func getCallGasEstimateNoop() GetCallGasEstimateFunc {
	return func(ep common.Address, op *userop.UserOperation) (uint64, error) {
		//lint:ignore ST1005 This needs to match the bundler test spec.
		return 0, errors.New("Missing/invalid userOpHash")
	}
}

// GetCallGasEstimateWithClient returns an implementation of GetCallGasEstimateFunc that relies on an eth
// client to fetch an estimate for callGas.
func GetCallGasEstimateWithClient(eth *ethclient.Client) GetCallGasEstimateFunc {
	return func(ep common.Address, op *userop.UserOperation) (uint64, error) {
		return gas.CallGasEstimate(eth, ep, op)
	}
}
