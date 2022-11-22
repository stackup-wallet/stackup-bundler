package relay

import (
	"math/big"
	"net/http"

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
	db     *badger.DB
	logger logr.Logger
}

// New initializes a new EOA relayer for sending batches to the EntryPoint with IP throttling protection.
func New(db *badger.DB, l logr.Logger) *Relayer {
	return &Relayer{
		db:     db,
		logger: l.WithName("relayer"),
	}
}

// FilterByClient is a custom Gin middleware used to prevent requests from banned clients from adding their
// userOps to the mempool. Identifiers are prioritized by the following values:
//  1. X-Bundler-Client-Id header: See UseClientIDHeader
//  2. X-Forwarded-By header: The first IP address in the array which is assumed to be the client
//  3. Request.RemoteAddr: The remote IP address
func (r *Relayer) FilterByClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := r.logger.WithName("filter_by_client")

		isBanned := false
		err := r.db.View(func(txn *badger.Txn) error {
			opsSeen, opsIncluded, err := getOpsCountByClientID(txn, ginutils.GetClientIPFromXFF(c))
			if err != nil {
				return err
			}

			OpsFailed := opsSeen - opsIncluded
			if OpsFailed < banThreshold {
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

		if !isBanned {
			c.Status(http.StatusForbidden)
			c.Abort()
		}
	}
}

// MapRequestIDToClientID is a custom Gin middleware used to map a userOp requestID to a client
// identifier (e.g. IP address).
func (r *Relayer) MapRequestIDToClientID(chainID *big.Int) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := r.logger.WithName("map_request_id_to_client_id")

		req, _ := c.Get("JsonRpcRequest")
		json := req.(map[string]any)
		if json["method"] != "eth_sendUserOperation" {
			return
		}

		params := json["params"].([]any)
		data := params[0].(map[string]any)
		ep := params[1].(string)
		op, err := userop.New(data)
		if err != nil {
			l.Error(err, "map_request_id_to_client_id failed")
			c.Status(http.StatusInternalServerError)
			return
		}

		rid := op.GetRequestID(common.HexToAddress(ep), chainID).String()
		cid := ginutils.GetClientIPFromXFF(c)
		err = r.db.Update(func(txn *badger.Txn) error {
			err := mapRequestIDToClientID(txn, rid, cid)
			if err != nil {
				return err
			}

			err = incrementOpsSeenByClientID(txn, cid)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			l.Error(err, "map_request_id_to_client_id failed")
			c.Status(http.StatusInternalServerError)
			return
		}
	}
}

// SendUserOperation returns a BatchHandler that accepts a batch and sends it as a regular EOA transaction.
// It can also map a userOp request ID to a Client ID (e.g. IP address) in order to mitigate DoS attacks.
func (r *Relayer) SendUserOperation(
	eoa *signer.EOA,
	eth *ethclient.Client,
	beneficiary common.Address,
) modules.BatchHandlerFunc {
	return func(ctx *modules.BatchHandlerCtx) error {
		err := r.db.Update(func(txn *badger.Txn) error {
			// Delete any request ID entries from dropped userOps.
			if len(ctx.PendingRemoval) > 0 {
				rids := getRequestIDsFromOps(ctx.EntryPoint, ctx.ChainID, ctx.PendingRemoval...)
				if err := removeRequestIDEntries(txn, rids...); err != nil {
					return err
				}
			}

			// Estimate gas for handleOps() and drop all userOps that cause unexpected reverts.
			var gas uint64
			for len(ctx.Batch) > 0 {
				est, revert, err := entrypoint.EstimateHandleOpsGas(
					eoa,
					eth,
					ctx.ChainID,
					ctx.EntryPoint,
					ctx.Batch,
					beneficiary,
				)

				if err != nil {
					return err
				} else if revert != nil {
					ctx.MarkOpIndexForRemoval(revert.OpIndex)

					rids := getRequestIDsFromOps(ctx.EntryPoint, ctx.ChainID, ctx.PendingRemoval...)
					if err := removeRequestIDEntries(txn, rids...); err != nil {
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
					eoa,
					eth,
					ctx.ChainID,
					ctx.EntryPoint,
					ctx.Batch,
					beneficiary,
					gas,
					ctx.Batch[0].MaxPriorityFeePerGas,
					ctx.Batch[0].MaxFeePerGas,
				)

				if err != nil {
					return err
				} else if revert != nil {
					ctx.MarkOpIndexForRemoval(revert.OpIndex)

					rids := getRequestIDsFromOps(ctx.EntryPoint, ctx.ChainID, ctx.PendingRemoval...)
					if err := removeRequestIDEntries(txn, rids...); err != nil {
						return err
					}
				} else {
					ctx.Data["txn_hash"] = t.Hash().String()
					break
				}
			}

			// Delete remaining request ID entries from submitted userOps.
			rids := getRequestIDsFromOps(ctx.EntryPoint, ctx.ChainID, ctx.Batch...)
			if err := incrementOpsIncludedByRequestIDs(txn, rids...); err != nil {
				return err
			}
			if err := removeRequestIDEntries(txn, rids...); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	}
}
