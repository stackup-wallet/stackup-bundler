package client

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
)

type GetUserOpReceiptFunc = func(hash string, ep common.Address) (*entrypoint.UserOperationReceipt, error)

func getUserOpReceiptNoop() GetUserOpReceiptFunc {
	return func(hash string, ep common.Address) (*entrypoint.UserOperationReceipt, error) {
		//lint:ignore ST1005 This needs to match the bundler test spec.
		return nil, errors.New("Missing/invalid userOpHash")
	}
}

func GetUserOpReceiptWithEthClient(eth *ethclient.Client) GetUserOpReceiptFunc {
	return func(hash string, ep common.Address) (*entrypoint.UserOperationReceipt, error) {
		return entrypoint.GetUserOperationReceipt(eth, hash, ep)
	}
}
