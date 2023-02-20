package transaction

import (
	bytesPkg "bytes"
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/reverts"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func toAbiType(batch []*userop.UserOperation) []entrypoint.UserOperation {
	ops := []entrypoint.UserOperation{}
	for _, op := range batch {
		ops = append(ops, entrypoint.UserOperation(*op))
	}

	return ops
}

// EstimateHandleOpsGas returns a gas estimate required to call handleOps() with a given batch. A failed call
// will return the cause of the revert.
func EstimateHandleOpsGas(
	eoa *signer.EOA,
	eth *ethclient.Client,
	chainID *big.Int,
	entryPoint common.Address,
	batch []*userop.UserOperation,
	beneficiary common.Address,
) (gas uint64, revert *reverts.FailedOpRevert, err error) {
	ep, err := entrypoint.NewEntrypoint(entryPoint, eth)
	if err != nil {
		return 0, nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(eoa.PrivateKey, chainID)
	if err != nil {
		return 0, nil, err
	}
	auth.GasLimit = math.MaxUint64
	auth.NoSend = true

	tx, err := ep.HandleOps(auth, toAbiType(batch), beneficiary)
	if err != nil {
		return 0, nil, err
	}

	est, err := eth.EstimateGas(context.Background(), ethereum.CallMsg{
		From:       eoa.Address,
		To:         tx.To(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		GasFeeCap:  tx.GasFeeCap(),
		GasTipCap:  tx.GasTipCap(),
		Value:      tx.Value(),
		Data:       tx.Data(),
		AccessList: tx.AccessList(),
	})
	if err != nil {
		revert, err := reverts.NewFailedOp(err)
		if err != nil {
			return 0, nil, err
		}
		return 0, revert, nil
	}

	return est, nil, nil
}

// HandleOps calls handleOps() on the EntryPoint with a given batch, gas limit, and tip. A failed call will
// return the cause of the revert.
func HandleOps(
	eoa *signer.EOA,
	eth *ethclient.Client,
	chainID *big.Int,
	entryPoint common.Address,
	batch []*userop.UserOperation,
	beneficiary common.Address,
	gas uint64,
) (txn *types.Transaction, revert *reverts.FailedOpRevert, err error) {
	ep, err := entrypoint.NewEntrypoint(entryPoint, eth)
	if err != nil {
		return nil, nil, err
	}
	tip, err := eth.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(eoa.PrivateKey, chainID)
	if err != nil {
		return nil, nil, err
	}
	auth.GasLimit = gas
	auth.GasTipCap = tip

	txn, err = ep.HandleOps(auth, toAbiType(batch), beneficiary)
	if err != nil {
		revert, err := reverts.NewFailedOp(err)
		if err != nil {
			return nil, nil, err
		}
		return nil, revert, nil
	}

	return txn, nil, nil
}

// CreateRawHandleOps returns a raw transaction string that calls handleOps() on the EntryPoint with a given
// batch, gas limit, and tip.
func CreateRawHandleOps(
	eoa *signer.EOA,
	eth *ethclient.Client,
	chainID *big.Int,
	entryPoint common.Address,
	batch []*userop.UserOperation,
	beneficiary common.Address,
	gas uint64,
	baseFee *big.Int,
) (string, error) {
	ep, err := entrypoint.NewEntrypoint(entryPoint, eth)
	if err != nil {
		return "", err
	}
	tip, err := eth.SuggestGasTipCap(context.Background())
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(eoa.PrivateKey, chainID)
	if err != nil {
		return "", err
	}
	auth.GasLimit = gas
	auth.GasTipCap = tip
	auth.GasFeeCap = big.NewInt(0).Add(baseFee, tip)
	auth.NoSend = true
	tx, err := ep.HandleOps(auth, toAbiType(batch), beneficiary)
	if err != nil {
		return "", err
	}

	ts := types.Transactions{tx}
	rawTxBytes := new(bytesPkg.Buffer)
	ts.EncodeIndex(0, rawTxBytes)
	return hexutil.Encode(rawTxBytes.Bytes()), nil
}
