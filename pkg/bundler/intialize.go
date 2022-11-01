package bundler

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
)

func New(ethClient *ethclient.Client, mempool *mempool.Interface, entryPoint common.Address) (*Instance, error) {
	return &Instance{
		ethClient:  ethClient,
		mempool:    mempool,
		entryPoint: entryPoint,
	}, nil
}
