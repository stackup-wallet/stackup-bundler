package entities

import (
	"github.com/ethereum/go-ethereum/common"
)

type ReputationOverride struct {
	Address     common.Address `json:"address"`
	OpsSeen     int            `json:"opsSeen"`
	OpsIncluded int            `json:"opsIncluded"`
}

// ReputationConstants are a collection of values for determining the appropriate status of a UserOperation
// coming into the mempool.
type ReputationConstants struct {
	MinUnstakeDelay                int
	MinStakeValue                  int64
	SameSenderMempoolCount         int
	SameUnstakedEntityMempoolCount int
	ThrottledEntityMempoolCount    int
	ThrottledEntityLiveBlocks      int
	ThrottledEntityBundleCount     int
	MinInclusionRateDenominator    int
	ThrottlingSlack                int
	BanSlack                       int
}
