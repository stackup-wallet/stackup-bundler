// Package relay implements a module for private bundlers to send batches to the EntryPoint through regular
// EOA transactions.
package relay

import (
	"math/big"
	"net/http"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/stackup-wallet/stackup-bundler/internal/ginutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Relayer provides a module that can relay batches with a regular EOA. Relaying batches to the EntryPoint
// through a regular transaction comes with several important notes:
//
//   - The bundler will NOT be operating as a block builder.
//   - This opens the bundler up to frontrunning.
//   - In a naive solution, attackers can send a valid op and frontrun the batch to make that op invalid.
//   - This invalidates the entire batch and the bundler will have to pay for the failed transaction.
//
// In this case, the mitigation strategy is to throttle the sender by a unique identifier or IP address.
// If a sender submits a UserOperation that causes the batch to revert, then its ID is banned from sending
// anymore ops to the mempool. This is optimistic in the sense that it will not prevent every case but will
// mitigate malicious senders spamming the mempool.
//
// This will only work in the case of a private mempool and will not work in the P2P case where ops are
// propagated through the network and it is impossible to trust a sender's identifier.
type Relayer struct {
	db              *badger.DB
	eoa             *signer.EOA
	eth             *ethclient.Client
	chainID         *big.Int
	beneficiary     common.Address
	logger          logr.Logger
	bannedThreshold int
}

// New initializes a new EOA relayer for sending batches to the EntryPoint with IP throttling protection.
func New(
	db *badger.DB,
	eoa *signer.EOA,
	eth *ethclient.Client,
	chainID *big.Int,
	beneficiary common.Address,
	l logr.Logger,
) *Relayer {
	return &Relayer{
		db:              db,
		eoa:             eoa,
		eth:             eth,
		chainID:         chainID,
		beneficiary:     beneficiary,
		logger:          l.WithName("relayer"),
		bannedThreshold: DefaultBanThreshold,
	}
}

// SetBannedThreshold sets the limit for how many ops can be seen from a client without being included in a
// batch before it is banned. Default value is 3. A value of 0 will effectively disable client banning, which
// is useful for debugging.
func (r *Relayer) SetBannedThreshold(limit int) {
	r.bannedThreshold = limit
}

// FilterByClientID is a custom Gin middleware used to prevent requests from banned clients from adding their
// userOps to the mempool. Identifiers are prioritized by the following values:
//  1. X-Forwarded-By header: The first IP address in the array which is assumed to be the client
//  2. Request.RemoteAddr: The remote IP address
//
// This should be the first middleware on the RPC path.
func (r *Relayer) FilterByClientID() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := r.logger.WithName("filter_by_client")

		isBanned := false
		cid := ginutils.GetClientIPFromXFF(c)
		err := r.db.View(func(txn *badger.Txn) error {
			opsSeen, opsIncluded, err := getOpsCountByClientID(txn, cid)
			if err != nil {
				return err
			}
			l = l.
				WithValues("client_id", cid).
				WithValues("opsSeen", opsSeen).
				WithValues("opsIncluded", opsIncluded)

			OpsFailed := opsSeen - opsIncluded
			if r.bannedThreshold == NoBanThreshold || OpsFailed < r.bannedThreshold {
				return nil
			}

			isBanned = true
			return nil
		})
		if err != nil {
			l.Error(err, "filter_by_client failed")
			c.Status(http.StatusInternalServerError)
			c.Abort()
		}

		if isBanned {
			l.Info("client banned")
			c.Status(http.StatusForbidden)
			c.Abort()
		} else {
			l.Info("client ok")
		}
	}
}

// MapUserOpHashToClientID is a custom Gin middleware used to map a userOpHash to a clientID. This
// should be placed after the main method call on the RPC path.
func (r *Relayer) MapUserOpHashToClientID() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := r.logger.WithName("map_userop_hash_to_client_id")

		req, _ := c.Get("json-rpc-request")
		json := req.(map[string]any)
		if json["method"] != "eth_sendUserOperation" {
			return
		}

		params := json["params"].([]any)
		data := params[0].(map[string]any)
		ep := params[1].(string)
		op, err := userop.New(data)
		if err != nil {
			l.Error(err, "map_userop_hash_to_client_id failed")
			c.Status(http.StatusInternalServerError)
			return
		}

		hash := op.GetUserOpHash(common.HexToAddress(ep), r.chainID).String()
		cid := ginutils.GetClientIPFromXFF(c)
		l = l.
			WithValues("userop_hash", hash).
			WithValues("client_id", cid)
		err = r.db.Update(func(txn *badger.Txn) error {
			err := mapUserOpHashToClientID(txn, hash, cid)
			if err != nil {
				return err
			}

			return incrementOpsSeenByClientID(txn, cid)
		})
		if err != nil {
			l.Error(err, "map_userop_hash_to_client_id failed")
			c.Status(http.StatusInternalServerError)
			return
		}
	}
}

// SendUserOperation returns a BatchHandler that is used by the Bundler to send batches in a regular EOA
// transaction. It uses the mapping of userOpHash to client ID created by the Gin middleware in order to
// mitigate DoS attacks.
func (r *Relayer) SendUserOperation() modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		// TODO: Increment badger nextTxnTs to read latest data from MapUserOpHashToClientID.
		time.Sleep(5 * time.Millisecond)

		var del []string
		err := r.db.Update(func(txn *badger.Txn) error {
			// Delete any userOpHash entries from dropped userOps.
			if len(ctx.PendingRemoval) > 0 {
				hashes := getUserOpHashesFromOps(ctx.EntryPoint, ctx.ChainID, ctx.PendingRemoval...)
				if err := removeUserOpHashEntries(txn, hashes...); err != nil {
					return err
				}
			}

			// Estimate gas for handleOps() and drop all userOps that cause unexpected reverts.
			var gas uint64
			for len(ctx.Batch) > 0 {
				est, revert, err := entrypoint.EstimateHandleOpsGas(
					r.eoa,
					r.eth,
					ctx.ChainID,
					ctx.EntryPoint,
					ctx.Batch,
					r.beneficiary,
				)

				if err != nil {
					return err
				} else if revert != nil {
					ctx.MarkOpIndexForRemoval(revert.OpIndex)

					hashes := getUserOpHashesFromOps(ctx.EntryPoint, ctx.ChainID, ctx.PendingRemoval...)
					if err := removeUserOpHashEntries(txn, hashes...); err != nil {
						return err
					}
				} else {
					gas = est
					break
				}
			}

			// Call handleOps() with gas estimate and drop all userOps that cause unexpected reverts.
			for len(ctx.Batch) > 0 {
				t, revert, err := entrypoint.HandleOps(
					r.eoa,
					r.eth,
					ctx.ChainID,
					ctx.EntryPoint,
					ctx.Batch,
					r.beneficiary,
					gas,
				)

				if err != nil {
					return err
				} else if revert != nil {
					ctx.MarkOpIndexForRemoval(revert.OpIndex)

					hashes := getUserOpHashesFromOps(ctx.EntryPoint, ctx.ChainID, ctx.PendingRemoval...)
					if err := removeUserOpHashEntries(txn, hashes...); err != nil {
						return err
					}
				} else {
					ctx.Data["txn_hash"] = t.Hash().String()
					break
				}
			}

			hashes := getUserOpHashesFromOps(ctx.EntryPoint, ctx.ChainID, ctx.Batch...)
			del = append([]string{}, hashes...)
			return incrementOpsIncludedByUserOpHashes(txn, hashes...)
		})
		if err != nil {
			return err
		}

		// Delete remaining userOpHash entries from submitted userOps.
		// Perform update in new txn to avoid db conflicts.
		err = r.db.Update(func(txn *badger.Txn) error {
			return removeUserOpHashEntries(txn, del...)
		})
		if err != nil {
			return err
		}

		return nil
	}
}
