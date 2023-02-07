// Package builder implements a module for bundlers to act as MEV searchers and send batches to the EntryPoint
// via a Block Builder API that supports eth_sendBundle.
package builder

import (
	mapset "github.com/deckarep/golang-set/v2"
)

var (
	// CompatibleChainIDs is a set of chainIDs that support the Block Builder API.
	CompatibleChainIDs = mapset.NewSet[uint64](1, 5)
)
