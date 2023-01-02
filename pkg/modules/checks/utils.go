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

// GetBaseFeeFunc provides a general interface for retrieving the closest estimate for basefee to allow for
// timely execution of a transaction.
type GetBaseFeeFunc = func() (*big.Int, error)

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

// getBaseFeeWithEthClient returns a GetBaseFeeFunc using an eth client.
func getBaseFeeWithEthClient(eth *ethclient.Client) GetBaseFeeFunc {
	return func() (*big.Int, error) {
		bn, err := eth.BlockNumber(context.Background())
		if err != nil {
			return nil, err
		}

		blk, _ := eth.BlockByNumber(context.Background(), big.NewInt(0).SetUint64(bn))
		if err != nil {
			return nil, err
		}
		return blk.BaseFee(), nil
	}
}
