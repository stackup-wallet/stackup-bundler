package transaction

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func wait(
	ctx context.Context,
	eth *ethclient.Client,
	tx *types.Transaction,
	d time.Duration,
) (*types.Receipt, error) {
	queryTicker := time.NewTicker(d)
	defer queryTicker.Stop()

	for {
		receipt, err := eth.TransactionReceipt(ctx, tx.Hash())
		if err == nil {
			return receipt, nil
		}

		// Wait for the next round.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}
