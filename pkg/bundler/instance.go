package bundler

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type BatchHandlerFunc func(batch []*userop.UserOperation) error
type ErrorHandlerFunc func(err error)

type Instance struct {
	ethClient            *ethclient.Client
	mempool              *mempool.Interface
	supportedEntryPoints []common.Address
	batchHandler         BatchHandlerFunc
	errorHandler         ErrorHandlerFunc
}

func (i *Instance) SetBatchHandlerFunc(handler BatchHandlerFunc) {
	i.batchHandler = handler
}

func (i *Instance) SetErrorHandlerFunc(handler ErrorHandlerFunc) {
	i.errorHandler = handler
}

func (i *Instance) Run() error {
	go func(i *Instance) {
		for {
			for _, ep := range i.supportedEntryPoints {
				batch, err := i.mempool.BundleOps(ep)
				if err != nil {
					i.errorHandler(err)
					continue
				}
				batch = filterSender(batch)
				batch = filterPaymaster(batch)

				err = i.batchHandler(batch)
				if err != nil {
					i.errorHandler(err)
					continue
				}

				err = i.mempool.RemoveOps(ep, getSenders(batch))
				if err != nil {
					i.errorHandler(err)
					continue
				}
			}

			time.Sleep(5 * time.Second)
		}
	}(i)

	return nil
}
