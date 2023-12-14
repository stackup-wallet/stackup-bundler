package builder

import (
	"errors"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/stackup-wallet/stackup-bundler/internal/config"
)

var (
	// CompatibleChainIDs is a set of chainIDs that support the Block Builder API.
	CompatibleChainIDs = mapset.NewSet(
		config.EthereumChainID.Uint64(),
		config.GoerliChainID.Uint64(),
		config.SepoliaChainID.Uint64(),
	)

	DefaultWaitTimeout = 72 * time.Second

	ErrFlashbotsBroadcastBundle = errors.New("flashbots broadcast bundle error")
)
