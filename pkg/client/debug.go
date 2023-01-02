package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/bundler"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
)

// Debug exposes methods used for testing the bundler. These should not be made available in production.
type Debug struct {
	eoa         *signer.EOA
	eth         *ethclient.Client
	mempool     *mempool.Mempool
	bundler     *bundler.Bundler
	chainID     *big.Int
	entrypoint  common.Address
	beneficiary common.Address
}

func NewDebug(
	eoa *signer.EOA,
	eth *ethclient.Client,
	mempool *mempool.Mempool,
	bundler *bundler.Bundler,
	chainID *big.Int,
	entrypoint common.Address,
	beneficiary common.Address,
) *Debug {
	return &Debug{eoa, eth, mempool, bundler, chainID, entrypoint, beneficiary}
}

// ClearState clears the bundler mempool and reputation data of paymasters/accounts/factories/aggregators.
func (d *Debug) ClearState() (string, error) {
	if err := d.mempool.Clear(); err != nil {
		return "", err
	}

	return "ok", nil
}

// DumpMempool dumps the current UserOperations mempool in order of arrival.
func (d *Debug) DumpMempool(ep string) ([]map[string]any, error) {
	ops, err := d.mempool.Dump(common.HexToAddress(ep))
	if err != nil {
		return []map[string]any{}, err
	}

	res := []map[string]any{}
	for _, op := range ops {
		data, err := op.MarshalJSON()
		if err != nil {
			return []map[string]any{}, err
		}

		item := make(map[string]any)
		if err := json.Unmarshal(data, &item); err != nil {
			return []map[string]any{}, err
		}

		res = append(res, item)
	}

	return res, nil
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
	)
	if err != nil {
		return "", err
	} else if revert != nil {
		return "", errors.New("debug: bad batch during call")
	}

	return txn.Hash().String(), nil
}

// SetBundlingMode allows the bundler to be stopped so that an explicit call to debug_bundler_sendBundleNow is
// required to send a bundle.
func (d *Debug) SetBundlingMode(mode string) (string, error) {
	switch mode {
	case "manual":
		d.bundler.Stop()
	case "auto":
		if err := d.bundler.Run(); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("debug: unrecognized mode %s", mode)
	}

	return "ok", nil
}
