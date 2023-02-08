package builder

import (
	mapset "github.com/deckarep/golang-set/v2"
)

var (
	// CompatibleChainIDs is a set of chainIDs that support the Block Builder API.
	CompatibleChainIDs = mapset.NewSet[uint64](1, 5)
)
