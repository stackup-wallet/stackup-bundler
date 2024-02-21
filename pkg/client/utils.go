package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/filter"
	"github.com/stackup-wallet/stackup-bundler/pkg/fees"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/state"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// GetUserOpReceiptFunc is a general interface for fetching a UserOperationReceipt given a userOpHash,
// EntryPoint address, and block range.
type GetUserOpReceiptFunc = func(hash string, ep common.Address, blkRange uint64) (*filter.UserOperationReceipt, error)

func getUserOpReceiptNoop() GetUserOpReceiptFunc {
	return func(hash string, ep common.Address, blkRange uint64) (*filter.UserOperationReceipt, error) {
		return nil, nil
	}
}

// GetUserOpReceiptWithEthClient returns an implementation of GetUserOpReceiptFunc that relies on an eth
// client to fetch a UserOperationReceipt.
func GetUserOpReceiptWithEthClient(eth *ethclient.Client) GetUserOpReceiptFunc {
	return func(hash string, ep common.Address, blkRange uint64) (*filter.UserOperationReceipt, error) {
		return filter.GetUserOperationReceipt(eth, hash, ep, blkRange)
	}
}

// GetGasPricesFunc is a general interface for fetching values for maxFeePerGas and maxPriorityFeePerGas.
type GetGasPricesFunc = func() (*fees.GasPrices, error)

func getGasPricesNoop() GetGasPricesFunc {
	return func() (*fees.GasPrices, error) {
		return &fees.GasPrices{
			MaxFeePerGas:         big.NewInt(0),
			MaxPriorityFeePerGas: big.NewInt(0),
		}, nil
	}
}

// GetGasPricesWithEthClient returns an implementation of GetGasPricesFunc that relies on an eth client to
// fetch values for maxFeePerGas and maxPriorityFeePerGas.
func GetGasPricesWithEthClient(eth *ethclient.Client) GetGasPricesFunc {
	return func() (*fees.GasPrices, error) {
		return fees.NewGasPrices(eth)
	}
}

// GetGasEstimateFunc is a general interface for fetching an estimate for verificationGasLimit and
// callGasLimit given a userOp and EntryPoint address.
type GetGasEstimateFunc = func(
	ep common.Address,
	op *userop.UserOperation,
	sos state.OverrideSet,
) (verificationGas uint64, callGas uint64, err error)

func getGasEstimateNoop() GetGasEstimateFunc {
	return func(
		ep common.Address,
		op *userop.UserOperation,
		sos state.OverrideSet,
	) (verificationGas uint64, callGas uint64, err error) {
		return 0, 0, nil
	}
}

// GetGasEstimateWithEthClient returns an implementation of GetGasEstimateFunc that relies on an eth client to
// fetch an estimate for verificationGasLimit and callGasLimit.
func GetGasEstimateWithEthClient(
	rpc *rpc.Client,
	ov *gas.Overhead,
	chain *big.Int,
	maxGasLimit *big.Int,
	tracer string,
) GetGasEstimateFunc {
	return func(
		ep common.Address,
		op *userop.UserOperation,
		sos state.OverrideSet,
	) (verificationGas uint64, callGas uint64, err error) {
		return gas.EstimateGas(&gas.EstimateInput{
			Rpc:         rpc,
			EntryPoint:  ep,
			Op:          op,
			Sos:         sos,
			Ov:          ov,
			ChainID:     chain,
			MaxGasLimit: maxGasLimit,
			Tracer:      tracer,
		})
	}
}

// GetUserOpByHashFunc is a general interface for fetching a UserOperation given a userOpHash, EntryPoint
// address, chain ID, and block range.
type GetUserOpByHashFunc func(hash string, ep common.Address, chain *big.Int, blkRange uint64) (*filter.HashLookupResult, error)

func getUserOpByHashNoop() GetUserOpByHashFunc {
	return func(hash string, ep common.Address, chain *big.Int, blkRange uint64) (*filter.HashLookupResult, error) {
		return nil, nil
	}
}

// GetUserOpByHashWithEthClient returns an implementation of GetUserOpByHashFunc that relies on an eth client
// to fetch a UserOperation.
func GetUserOpByHashWithEthClient(eth *ethclient.Client) GetUserOpByHashFunc {
	return func(hash string, ep common.Address, chain *big.Int, blkRange uint64) (*filter.HashLookupResult, error) {
		return filter.GetUserOperationByHash(eth, hash, ep, chain, blkRange)
	}
}
