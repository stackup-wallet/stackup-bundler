package relay

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/noop"
)

// New initializes a new EOA relayer for sending batches to the EntryPoint with IP throttling protection.
func New(db *badger.DB) (*Relayer, error) {
	return &Relayer{
		db:                    db,
		errorHandler:          noop.ErrorHandler,
		clientIDHeaderEnabled: false,
	}, nil
}
