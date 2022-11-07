package relay

import (
	"math/big"

	"github.com/dgraph-io/badger/v3"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/noop"
)

// New initializes a new EOA relayer for sending batches to the EntryPoint with IP throttling protection.
func New(db *badger.DB, chainID *big.Int) (*Relayer, error) {
	return &Relayer{
		db:                    db,
		chainID:               chainID,
		errorHandler:          noop.ErrorHandler,
		clientIDHeaderEnabled: false,
	}, nil
}
