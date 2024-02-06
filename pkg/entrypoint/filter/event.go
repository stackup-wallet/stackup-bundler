package filter

import (
	"context"
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
	blkRange uint64,
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
	subBlkRange := big.NewInt(0).Sub(toBlk, big.NewInt(0).SetUint64(blkRange))
	if subBlkRange.Cmp(startBlk) > 0 {
		startBlk = subBlkRange
	}

	return ep.FilterUserOperationEvent(
		&bind.FilterOpts{Start: startBlk.Uint64()},
		[][32]byte{common.HexToHash(userOpHash)},
		[]common.Address{},
		[]common.Address{},
	)
}
