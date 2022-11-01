package bundler

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
)

type Instance struct {
	ethClient  *ethclient.Client
	mempool    *mempool.Interface
	entryPoint common.Address
}

func (i *Instance) ProcessBatch() (bool, error) {
	return true, nil
}
