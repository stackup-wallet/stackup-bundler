package paymaster

import (
	"errors"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
)

// Reputation provides client and bundler modules to track the status of every Paymaster seen in a
// UserOperation.
type Reputation struct {
	db *badger.DB
}

// New returns an instance of a Reputation module to track paymaster status.
func New(db *badger.DB) *Reputation {
	return &Reputation{db}
}

// StatusCheck returns a UserOpHandler to determine if the userOp is allowed based on the paymaster status.
//  1. ok: Paymasters is allowed
//  2. throttled: No new ops from the Paymaster is allowed if one already exists. And it can only stays in
//     the pool for 10 blocks
//  3. banned: No ops from the Paymaster is allowed
func (r *Reputation) StatusCheck() modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		return r.db.Update(func(txn *badger.Txn) error {
			paymaster := ctx.UserOp.GetPaymaster()
			if paymaster == common.HexToAddress("0x") {
				return nil
			}

			status, err := getStatus(txn, paymaster)
			if err != nil {
				return err
			}

			if status == banned {
				return errors.New("paymaster: is currently banned")
			}
			// TODO: Implement logic for throttled status

			return nil
		})
	}
}

// IncOpsSeen returns a UserOpHandler that checks if a userOp has a paymaster and increments its opsSeen
// counter.
func (r *Reputation) IncOpsSeen() modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		return r.db.Update(func(txn *badger.Txn) error {
			paymaster := ctx.UserOp.GetPaymaster()
			if paymaster == common.HexToAddress("0x") {
				return nil
			}

			return incrementOpsSeenByPaymaster(txn, paymaster)
		})
	}
}

// IncOpsIncluded returns a BatchHandler that increments opsIncluded for paymasters with userOps in the
// batch. This module should be used last once batches have been sent.
func (r *Reputation) IncOpsIncluded() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		return r.db.Update(func(txn *badger.Txn) error {
			ps := mapset.NewSet[common.Address]()
			for _, op := range ctx.Batch {
				paymaster := op.GetPaymaster()
				if paymaster != common.HexToAddress("0x") {
					ps.Add(paymaster)
				}
			}

			return incrementOpsIncludedByPaymasters(txn, ps.ToSlice()...)
		})
	}
}
