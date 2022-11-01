package client

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
)

// Initializes a new ERC-4337 client with an ethClient instance
// and an array of supported EntryPoint addresses.
// The first address in the array is the preferred EntryPoint.
func New(ethClient *ethclient.Client, mempool *mempool.Interface, supportedEntryPoints []common.Address) (*Instance, error) {
	cid, err := ethClient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	return &Instance{
		ethClient:            ethClient,
		mempool:              mempool,
		chainID:              cid,
		supportedEntryPoints: supportedEntryPoints,
	}, nil
}
