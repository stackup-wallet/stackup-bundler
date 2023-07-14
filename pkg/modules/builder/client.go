// Package builder implements a module for bundlers to act as MEV searchers and send batches to the EntryPoint
// via a Block Builder API that supports eth_sendBundle.
package builder

import (
	"context"
	"errors"
	"fmt"
	"math/big"

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
	rpc               *flashbotsrpc.FlashbotsRPC
	beneficiary       common.Address
	blocksInTheFuture int
}

// New returns an instance of a BuilderClient with modules to send UserOperation bundles via the mev-boost
// process.
func New(
	eoa *signer.EOA,
	eth *ethclient.Client,
	fb *flashbotsrpc.FlashbotsRPC,
	beneficiary common.Address,
	blocksInTheFuture int,
) *BuilderClient {
	return &BuilderClient{
		eoa:               eoa,
		eth:               eth,
		rpc:               fb,
		beneficiary:       beneficiary,
		blocksInTheFuture: blocksInTheFuture,
	}
}

// SendUserOperation returns a BatchHandler that is used by the Bundler to send batches to a block builder
// that supports eth_callBundle and eth_sendBundle.
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
			GasPrice:    ctx.GasPrice,
			GasLimit:    0,
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
		blkNum := big.NewInt(0).SetUint64(bn)
		NxtBlkNum := big.NewInt(0).Add(blkNum, big.NewInt(1))
		mbf := ctx.BaseFee
		for i := 0; i < b.blocksInTheFuture; i++ {
			a := big.NewInt(0).Mul(mbf, big.NewInt(1125))
			b := big.NewInt(0).Div(a, big.NewInt(1000))
			mbf = big.NewInt(0).Add(b, big.NewInt(1))
		}
		opts.BaseFee = mbf

		// Call CreateRawHandleOps() with gas estimate and max base fee.
		rawTx, err := transaction.CreateRawHandleOps(&opts)
		if err != nil {
			return err
		}

		// Simulate bundle.
		callBundleArgs := flashbotsrpc.FlashbotsCallBundleParam{
			Txs:              []string{rawTx},
			BlockNumber:      hexutil.EncodeBig(NxtBlkNum),
			StateBlockNumber: "latest",
		}
		sim, err := b.rpc.FlashbotsCallBundle(b.eoa.PrivateKey, callBundleArgs)
		if err != nil {
			return err
		}
		if len(sim.Results) != 1 {
			return fmt.Errorf("unexpected simulation result length, want 1, got %d", len(sim.Results))
		}
		if sim.Results[0].Error != "" {
			// TODO: Implement better error handling and retry.
			return errors.New(sim.Results[0].Error)
		}

		// Send bundle for all blocks up to a future block number.
		// Note: Do not try to access bundleHash from results. Flashbots builder does not return it.
		for i := 0; i < b.blocksInTheFuture; i++ {
			futureBlkNum := big.NewInt(0).Add(blkNum, big.NewInt(int64(i)))
			sendBundleArgs := flashbotsrpc.FlashbotsSendBundleRequest{
				Txs:         []string{rawTx},
				BlockNumber: hexutil.EncodeBig(futureBlkNum),
			}
			_, err := b.rpc.FlashbotsSendBundle(b.eoa.PrivateKey, sendBundleArgs)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
