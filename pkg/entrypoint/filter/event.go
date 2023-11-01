package filter

import (
	"context"
	"github.com/spf13/viper"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
)

func filterUserOperationEvent(
	eth *ethclient.Client,
	userOpHash string,
	entryPoint common.Address,
) (*entrypoint.EntrypointUserOperationEventIterator, error) {
	ep, err := entrypoint.NewEntrypoint(entryPoint, eth)
	if err != nil {
		return nil, err
	}
	bn, err := eth.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	toBlk := big.NewInt(0).SetUint64(bn)
	startBlk := big.NewInt(0)
	getLogsStepSize := viper.GetInt64("erc4337_bundler_get_logs_step_size")
	subStepSizeBlk := big.NewInt(0).Sub(toBlk, big.NewInt(getLogsStepSize))
	if subStepSizeBlk.Cmp(startBlk) > 0 {
		startBlk = subStepSizeBlk
	}

	return ep.FilterUserOperationEvent(
		&bind.FilterOpts{Start: startBlk.Uint64()},
		[][32]byte{common.HexToHash(userOpHash)},
		[]common.Address{},
		[]common.Address{},
	)
}
