package bundler

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func New(ethClient *ethclient.Client, mempool *mempool.Interface, supportedEntryPoints []common.Address) (*Instance, error) {
	return &Instance{
		ethClient:            ethClient,
		mempool:              mempool,
		supportedEntryPoints: supportedEntryPoints,
		batchHandler:         func(batch []*userop.UserOperation) error { return nil },
		errorHandler:         func(err error) {},
	}, nil
}
