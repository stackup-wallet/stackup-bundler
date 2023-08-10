package expire

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
)

type ExpireHandler struct {
	seenAt map[common.Hash]time.Time
	ttl    time.Duration
}

// New returns an ExpireHandler which contains a BatchHandlerFunc to track and drop UserOperations that have
// been in the mempool for longer than the TTL duration.
func New(ttl time.Duration) *ExpireHandler {
	return &ExpireHandler{
		seenAt: make(map[common.Hash]time.Time),
		ttl:    ttl,
	}
}

// DropExpired returns a BatchHandlerFunc that will drop UserOperations from the mempool if it has been around
// for longer than the TTL duration.
func (e *ExpireHandler) DropExpired() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		end := len(ctx.Batch) - 1
		for i := end; i >= 0; i-- {
			hash := ctx.Batch[i].GetUserOpHash(ctx.EntryPoint, ctx.ChainID)
			if seenAt, ok := e.seenAt[hash]; !ok {
				e.seenAt[hash] = time.Now()
			} else if seenAt.Add(e.ttl).Before(time.Now()) {
				ctx.MarkOpIndexForRemoval(i)
			}
		}
		return nil
	}
}
