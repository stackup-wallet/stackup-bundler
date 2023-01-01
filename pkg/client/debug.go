package client

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
)

// Debug exposes methods used for testing the bundler. These should not be made available in production.
type Debug struct {
	eoa         *signer.EOA
	eth         *ethclient.Client
	mempool     *mempool.Mempool
	chainID     *big.Int
	entrypoint  common.Address
	beneficiary common.Address
}

func NewDebug(
	eoa *signer.EOA,
	eth *ethclient.Client,
	mempool *mempool.Mempool,
	chainID *big.Int,
	entrypoint common.Address,
	beneficiary common.Address,
) *Debug {
	return &Debug{eoa, eth, mempool, chainID, entrypoint, beneficiary}
}

// SendBundleNow forces the bundler to build and execute a bundle from the mempool as handleOps() transaction.
func (d *Debug) SendBundleNow() (string, error) {
	batch, err := d.mempool.BundleOps(d.entrypoint)
	if err != nil {
		return "", err
	}

	est, revert, err := entrypoint.EstimateHandleOpsGas(
		d.eoa,
		d.eth,
		d.chainID,
		d.entrypoint,
		batch,
		d.beneficiary,
	)
	if err != nil {
		return "", err
	} else if revert != nil {
		return "", errors.New("debug: bad batch during estimate")
	}

	txn, revert, err := entrypoint.HandleOps(
		d.eoa,
		d.eth,
		d.chainID,
		d.entrypoint,
		batch,
		d.beneficiary,
		est,
		big.NewInt((int64(est))),
		big.NewInt((int64(est))),
	)
	if err != nil {
		return "", err
	} else if revert != nil {
		return "", errors.New("debug: bad batch during call")
	}

	return txn.Hash().String(), nil
}
