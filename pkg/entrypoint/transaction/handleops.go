package transaction

import (
	bytesPkg "bytes"
	"context"
	"errors"
	"math"
	"math/big"
	"time"

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

// Opts contains all the fields required for submitting a transaction to call HandleOps on the EntryPoint
// contract.
type Opts struct {
	// Options for the network
	EOA     *signer.EOA
	Eth     *ethclient.Client
	ChainID *big.Int

	// Options for the EntryPoint
	EntryPoint  common.Address
	Batch       []*userop.UserOperation
	Beneficiary common.Address

	// Options for the EOA transaction
	BaseFee     *big.Int
	Tip         *big.Int
	GasPrice    *big.Int
	GasLimit    uint64
	WaitTimeout time.Duration
}

func toAbiType(batch []*userop.UserOperation) []entrypoint.UserOperation {
	ops := []entrypoint.UserOperation{}
	for _, op := range batch {
		ops = append(ops, entrypoint.UserOperation(*op))
	}

	return ops
}

// EstimateHandleOpsGas returns a gas estimate required to call handleOps() with a given batch. A failed call
// will return the cause of the revert.
func EstimateHandleOpsGas(opts *Opts) (gas uint64, revert *reverts.FailedOpRevert, err error) {
	ep, err := entrypoint.NewEntrypoint(opts.EntryPoint, opts.Eth)
	if err != nil {
		return 0, nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(opts.EOA.PrivateKey, opts.ChainID)
	if err != nil {
		return 0, nil, err
	}
	auth.GasLimit = math.MaxUint64
	auth.NoSend = true

	tx, err := ep.HandleOps(auth, toAbiType(opts.Batch), opts.Beneficiary)
	if err != nil {
		return 0, nil, err
	}

	est, err := opts.Eth.EstimateGas(context.Background(), ethereum.CallMsg{
		From:       opts.EOA.Address,
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

// HandleOps submits a transaction to send a batch of UserOperations to the EntryPoint.
func HandleOps(opts *Opts) (txn *types.Transaction, err error) {
	ep, err := entrypoint.NewEntrypoint(opts.EntryPoint, opts.Eth)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(opts.EOA.PrivateKey, opts.ChainID)
	if err != nil {
		return nil, err
	}
	auth.GasLimit = opts.GasLimit

	nonce, err := opts.Eth.NonceAt(context.Background(), opts.EOA.Address, nil)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	if opts.BaseFee != nil && opts.Tip != nil {
		auth.GasTipCap = SuggestMeanGasTipCap(opts.Tip, opts.Batch)
		auth.GasFeeCap = SuggestMeanGasFeeCap(opts.BaseFee, opts.Tip, opts.Batch)
	} else if opts.GasPrice != nil {
		auth.GasPrice = SuggestMeanGasPrice(opts.GasPrice, opts.Batch)
	} else {
		return nil, errors.New("transaction: either the dynamic or legacy gas fees must be set")
	}

	txn, err = ep.HandleOps(auth, toAbiType(opts.Batch), opts.Beneficiary)
	if err != nil {
		return nil, err
	} else if opts.WaitTimeout == 0 {
		// Don't wait for transaction to be included. All userOps in the current batch will be dropped
		// regardless of the transaction status.
		return txn, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.WaitTimeout)
	defer cancel()
	if receipt, err := bind.WaitMined(ctx, opts.Eth, txn); err != nil {
		return nil, err
	} else if receipt.Status == types.ReceiptStatusFailed {
		// Return an error here so that the current batch stays in the mempool. In the next bundler iteration,
		// the offending userOps will be dropped during gas estimation.
		return nil, errors.New("transaction: failed status")
	}
	return txn, nil
}

// CreateRawHandleOps returns a raw transaction string that calls handleOps() on the EntryPoint with a given
// batch, gas limit, and tip.
func CreateRawHandleOps(opts *Opts) (string, error) {
	ep, err := entrypoint.NewEntrypoint(opts.EntryPoint, opts.Eth)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(opts.EOA.PrivateKey, opts.ChainID)
	if err != nil {
		return "", err
	}
	auth.GasLimit = opts.GasLimit
	auth.NoSend = true
	if opts.BaseFee != nil {
		tip, err := opts.Eth.SuggestGasTipCap(context.Background())
		if err != nil {
			return "", err
		}

		auth.GasTipCap = tip
		auth.GasFeeCap = big.NewInt(0).Add(opts.BaseFee, tip)
	}

	tx, err := ep.HandleOps(auth, toAbiType(opts.Batch), opts.Beneficiary)
	if err != nil {
		return "", err
	}

	ts := types.Transactions{tx}
	rawTxBytes := new(bytesPkg.Buffer)
	ts.EncodeIndex(0, rawTxBytes)
	return hexutil.Encode(rawTxBytes.Bytes()), nil
}
