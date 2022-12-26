package checks

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
)

// GetCodeFunc provides a general interface for retrieving the bytecode for a given address.
type GetCodeFunc = func(addr common.Address) ([]byte, error)

// GetStakeFunc provides a general interface for retrieving the EntryPoint stake for a given address.
type GetStakeFunc = func(entity common.Address) (*entrypoint.IStakeManagerDepositInfo, error)

// GetGasTipFunc provides a general interface for retrieving an EIP-1559 style gas tip to allow for timely
// execution of a transaction.
type GetGasTipFunc = func() (*big.Int, error)

// getCodeWithEthClient returns a GetCodeFunc that uses an eth client to call eth_getCode.
func getCodeWithEthClient(eth *ethclient.Client) GetCodeFunc {
	return func(addr common.Address) ([]byte, error) {
		return eth.CodeAt(context.Background(), addr, nil)
	}
}

// getStakeWithEthClient returns a GetStakeFunc that uses an EntryPoint binding to get stake info and adds it
// to the current context.
func getStakeWithEthClient(ctx *modules.UserOpHandlerCtx, eth *ethclient.Client) (GetStakeFunc, error) {
	ep, err := entrypoint.NewEntrypoint(ctx.EntryPoint, eth)
	if err != nil {
		return nil, err
	}

	return func(addr common.Address) (*entrypoint.IStakeManagerDepositInfo, error) {
		dep, err := ep.GetDepositInfo(nil, addr)
		if err != nil {
			return nil, err
		}

		ctx.AddDepositInfo(addr, &dep)
		return &dep, nil
	}, nil
}

// getGasTipWithEthClient returns a GetGasTipFunc that uses an eth client to call eth_maxPriorityFeePerGas.
func getGasTipWithEthClient(eth *ethclient.Client) GetGasTipFunc {
	return func() (*big.Int, error) {
		return eth.SuggestGasTipCap(context.Background())
	}
}
