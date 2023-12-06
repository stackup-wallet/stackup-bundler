package transaction

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ToRawTxHex Takes a Geth types.Transaction and returns the encoded raw hex string.
func ToRawTxHex(txn *types.Transaction) string {
	rawTxn := new(bytes.Buffer)
	types.Transactions{txn}.EncodeIndex(0, rawTxn)
	return hexutil.Encode(rawTxn.Bytes())
}

// Wait blocks the process until a given transaction has been included on-chain or timeout has been reached.
func Wait(txn *types.Transaction, eth *ethclient.Client, timeout time.Duration) (*types.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if receipt, err := bind.WaitMined(ctx, eth, txn); err != nil {
		return nil, err
	} else if receipt.Status == types.ReceiptStatusFailed {
		// Return an error here so that the current batch stays in the mempool. In the next bundler iteration,
		// the offending userOps will be dropped during gas estimation.
		return nil, errors.New("transaction: failed status")
	}
	return txn, nil
}
