// Package entities implements modules for reputation scoring and throttling/banning of entities as specified
// in EIP-4337.
package entities

import (
	stdErr "errors"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
)

// Reputation provides Client and Bundler modules to track the reputation of every entity seen in a
// UserOperation.
type Reputation struct {
	db       *badger.DB
	eth      *ethclient.Client
	repConst *ReputationConstants
}

// New returns an instance of a Reputation object to track and appropriately process userOps by entity status.
func New(db *badger.DB, eth *ethclient.Client, repConst *ReputationConstants) *Reputation {
	return &Reputation{db, eth, repConst}
}

// CheckStatus returns a UserOpHandler that is used by the Client to determine if the userOp is allowed based
// on the entities status.
//  1. ok: entity is allowed
//  2. throttled: No new ops from the entity is allowed if one already exists. And it can only stays in
//     the pool for 10 blocks
//  3. banned: No ops from the entity is allowed
func (r *Reputation) CheckStatus() modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		return r.db.Update(func(txn *badger.Txn) error {
			if status, err := getStatus(txn, ctx.UserOp.Sender, r.repConst); err != nil {
				return err
			} else if status == banned {
				return errors.NewRPCError(
					errors.BANNED_OR_THROTTLED_ENTITY,
					fmt.Sprintf("banned entity: %s", ctx.UserOp.Sender.Hex()),
					nil,
				)
			} else if status == throttled && len(ctx.GetPendingSenderOps()) == r.repConst.ThrottledEntityMempoolCount {
				return errors.NewRPCError(
					errors.BANNED_OR_THROTTLED_ENTITY,
					fmt.Sprintf("throttled entity: %s", ctx.UserOp.Sender.Hex()),
					nil,
				)
			}

			factory := ctx.UserOp.GetFactory()
			if factory != common.HexToAddress("0x") {
				if status, err := getStatus(txn, factory, r.repConst); err != nil {
					return err
				} else if status == banned {
					return errors.NewRPCError(
						errors.BANNED_OR_THROTTLED_ENTITY,
						fmt.Sprintf("banned entity: %s", factory.Hex()),
						nil,
					)
				} else if status == throttled && len(ctx.GetPendingFactoryOps()) == r.repConst.ThrottledEntityMempoolCount {
					return errors.NewRPCError(
						errors.BANNED_OR_THROTTLED_ENTITY,
						fmt.Sprintf("throttled entity: %s", factory.Hex()),
						nil,
					)
				}
			}

			paymaster := ctx.UserOp.GetPaymaster()
			if paymaster != common.HexToAddress("0x") {
				if status, err := getStatus(txn, paymaster, r.repConst); err != nil {
					return err
				} else if status == banned {
					return errors.NewRPCError(
						errors.BANNED_OR_THROTTLED_ENTITY,
						fmt.Sprintf("banned entity: %s", paymaster.Hex()),
						nil,
					)
				} else if status == throttled && len(ctx.GetPendingPaymasterOps()) == r.repConst.ThrottledEntityMempoolCount {
					return errors.NewRPCError(
						errors.BANNED_OR_THROTTLED_ENTITY,
						fmt.Sprintf("throttled entity: %s", paymaster.Hex()),
						nil,
					)
				}
			}

			return nil
		})
	}
}

// ValidateOpLimit returns a UserOpHandler that is used by the Client to determine if the userOp is allowed
// based on the entities stake and the number of pending ops in the mempool.
func (r *Reputation) ValidateOpLimit() modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		pso := ctx.GetPendingSenderOps()
		sd := ctx.GetSenderDepositInfo()
		if !sd.Staked && len(pso) == r.repConst.SameSenderMempoolCount {
			return errors.NewRPCError(
				errors.INVALID_ENTITY_STAKE,
				fmt.Sprintf(
					"unstaked entity: %s exceeds pending ops limit of %d",
					ctx.UserOp.Sender.Hex(),
					r.repConst.SameSenderMempoolCount,
				),
				nil,
			)
		}

		factory := ctx.UserOp.GetFactory()
		if factory != common.HexToAddress("0x") {
			pfo := ctx.GetPendingFactoryOps()
			fd := ctx.GetFactoryDepositInfo()
			if !fd.Staked && len(pfo) == r.repConst.SameUnstakedEntityMempoolCount {
				return errors.NewRPCError(
					errors.INVALID_ENTITY_STAKE,
					fmt.Sprintf(
						"unstaked entity: %s exceeds pending ops limit of %d",
						factory.Hex(),
						r.repConst.SameUnstakedEntityMempoolCount,
					),
					nil,
				)
			}
		}

		paymaster := ctx.UserOp.GetPaymaster()
		if paymaster != common.HexToAddress("0x") {
			ppo := ctx.GetPendingPaymasterOps()
			pd := ctx.GetPaymasterDepositInfo()
			if !pd.Staked && len(ppo) == r.repConst.SameUnstakedEntityMempoolCount {
				return errors.NewRPCError(
					errors.INVALID_ENTITY_STAKE,
					fmt.Sprintf(
						"unstaked entity: %s exceeds pending ops limit of %d",
						paymaster.Hex(),
						r.repConst.SameUnstakedEntityMempoolCount,
					),
					nil,
				)
			}
		}

		return nil
	}
}

// IncOpsSeen returns a UserOpHandler that is used by the Client to increment the opsSeen counter for all
// included entities.
func (r *Reputation) IncOpsSeen() modules.UserOpHandlerFunc {
	return func(ctx *modules.UserOpHandlerCtx) error {
		return r.db.Update(func(txn *badger.Txn) error {
			var err error
			err = stdErr.Join(err, incrementOpsSeenByEntity(txn, ctx.UserOp.Sender))

			factory := ctx.UserOp.GetFactory()
			if factory != common.HexToAddress("0x") {
				err = stdErr.Join(err, incrementOpsSeenByEntity(txn, factory))
			}

			paymaster := ctx.UserOp.GetPaymaster()
			if paymaster != common.HexToAddress("0x") {
				err = stdErr.Join(err, incrementOpsSeenByEntity(txn, paymaster))
			}

			return err
		})
	}
}

// IncOpsIncluded returns a BatchHandler used by the Bundler to increment opsIncluded counters for all
// relevant entities in the batch. This module should be used last once batches have been sent.
func (r *Reputation) IncOpsIncluded() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		return r.db.Update(func(txn *badger.Txn) error {
			c := make(addressCounter)
			for _, op := range ctx.Batch {
				if _, ok := c[op.Sender]; !ok {
					c[op.Sender] = 0
				}
				c[op.Sender]++

				factory := op.GetFactory()
				if factory != common.HexToAddress("0x") {
					if _, ok := c[factory]; !ok {
						c[factory] = 0
					}

					c[factory]++
				}

				paymaster := op.GetPaymaster()
				if paymaster != common.HexToAddress("0x") {
					if _, ok := c[paymaster]; !ok {
						c[paymaster] = 0
					}

					c[paymaster]++
				}
			}

			return incrementOpsIncludedByEntity(txn, c)
		})
	}
}

func (r *Reputation) Override(entries []*ReputationOverride) error {
	return r.db.Update(func(txn *badger.Txn) error {
		var err error
		for _, entry := range entries {
			stdErr.Join(err, overrideEntity(txn, entry))
		}
		return err
	})
}
