// Package builder implements a module for bundlers to act as MEV searchers and send batches to the EntryPoint
// via a Block Builder API that supports eth_sendBundle.
package builder

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metachris/flashbotsrpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/transaction"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
)

// BuilderClient provides a connection to a block builder API to enable UserOperations to be sent through the
// mev-boost process.
type BuilderClient struct {
	eoa               *signer.EOA
	eth               *ethclient.Client
	rpc               *flashbotsrpc.BuilderBroadcastRPC
	beneficiary       common.Address
	blocksInTheFuture int
	waitTimeout       time.Duration
}

// New returns an instance of a BuilderClient with modules to send UserOperation bundles via the mev-boost
// process.
func New(
	eoa *signer.EOA,
	eth *ethclient.Client,
	fb *flashbotsrpc.BuilderBroadcastRPC,
	beneficiary common.Address,
	blocksInTheFuture int,
) *BuilderClient {
	return &BuilderClient{
		eoa:               eoa,
		eth:               eth,
		rpc:               fb,
		beneficiary:       beneficiary,
		blocksInTheFuture: blocksInTheFuture,
		waitTimeout:       DefaultWaitTimeout,
	}
}

// SetWaitTimeout sets the total time to wait for a transaction to be included. When a timeout is reached, the
// BatchHandler will throw an error if the transaction has not been included or has been included but with a
// failed status.
//
// The default value is 72 seconds. Setting the value to 0 will skip waiting for a transaction to be included.
func (b *BuilderClient) SetWaitTimeout(timeout time.Duration) {
	b.waitTimeout = timeout
}

// SendUserOperation returns a BatchHandler that is used by the Bundler to send batches to a block builder
// that supports eth_sendBundle.
func (b *BuilderClient) SendUserOperation() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		// Estimate gas for handleOps() and drop all userOps that cause unexpected reverts.
		opts := transaction.Opts{
			EOA:         b.eoa,
			Eth:         b.eth,
			ChainID:     ctx.ChainID,
			EntryPoint:  ctx.EntryPoint,
			Batch:       ctx.Batch,
			Beneficiary: b.beneficiary,
			BaseFee:     ctx.BaseFee,
			Tip:         ctx.Tip,
			GasPrice:    ctx.GasPrice,
			GasLimit:    0,
			NoSend:      true,
			WaitTimeout: b.waitTimeout,
		}
		for len(ctx.Batch) > 0 {
			est, revert, err := transaction.EstimateHandleOpsGas(&opts)

			if err != nil {
				return err
			} else if revert != nil {
				ctx.MarkOpIndexForRemoval(revert.OpIndex)
			} else {
				opts.GasLimit = est
				break
			}
		}

		// Calculate the max base fee up to a future block number.
		bn, err := b.eth.BlockNumber(context.Background())
		if err != nil {
			return err
		}
		nbn := big.NewInt(0).Add(big.NewInt(0).SetUint64(bn), big.NewInt(1))
		mbf := ctx.BaseFee
		for i := 0; i < b.blocksInTheFuture; i++ {
			a := big.NewInt(0).Mul(mbf, big.NewInt(1125))
			b := big.NewInt(0).Div(a, big.NewInt(1000))
			mbf = big.NewInt(0).Add(b, big.NewInt(1))
		}
		opts.BaseFee = mbf

		// Create no send transaction to the EntryPoint
		txn, err := transaction.HandleOps(&opts)
		if err != nil {
			return err
		}

		// Broadcast bundle to a list of ethereum block builders for all blocks up to a future block.
		shouldFail := true
		var errs error
		for i := 0; i < b.blocksInTheFuture; i++ {
			fbn := big.NewInt(0).Add(nbn, big.NewInt(int64(i)))
			sendBundleArgs := flashbotsrpc.FlashbotsSendBundleRequest{
				Txs:         []string{transaction.ToRawTxHex(txn)},
				BlockNumber: hexutil.EncodeBig(fbn),
			}

			results := b.rpc.BroadcastBundle(b.eoa.PrivateKey, sendBundleArgs)
			for _, result := range results {
				if result.Err != nil {
					errs = errors.Join(errs, result.Err)
				} else {
					shouldFail = false
				}
			}
		}

		// If there are no successful broadcast, return an error.
		if shouldFail {
			return fmt.Errorf("%w: \n\n%w", ErrFlashbotsBroadcastBundle, errs)
		}

		// Wait for transaction to be included on-chain.
		_, err = transaction.Wait(txn, opts.Eth, opts.WaitTimeout)
		return err
	}
}
