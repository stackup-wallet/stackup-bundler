package relay

import (
	"math/big"
	"net"
	"net/http"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
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
	db                    *badger.DB
	chainID               *big.Int
	errorHandler          modules.ErrorHandlerFunc
	clientIDHeaderEnabled bool
}

func (r *Relayer) getClientID(c *gin.Context) string {
	if r.clientIDHeaderEnabled && c.Request.Header.Get("x-bundler-client-id") != "" {
		return c.Request.Header.Get("x-bundler-client-id")
	}

	forwardHeader := c.Request.Header.Get("x-forwarded-for")
	firstAddress := strings.Split(forwardHeader, ",")[0]
	if net.ParseIP(strings.TrimSpace(firstAddress)) != nil {
		return firstAddress
	}

	return c.ClientIP()
}

// UseClientIDHeader allows bundlers to identify clients using any ID set in the X-Bundler-Client-Id header.
// This should only be turned on if incoming requests are from trusted sources.
func (r *Relayer) UseClientIDHeader(flag bool) {
	r.clientIDHeaderEnabled = flag
}

// SetErrorHandlerFunc defines a method for handling errors at any point of the process.
func (r *Relayer) SetErrorHandlerFunc(handler modules.ErrorHandlerFunc) {
	r.errorHandler = handler
}

// FilterByClient is a custom Gin middleware used to prevent requests from banned clients from adding their
// userOps to the mempool. Identifiers are prioritized by the following values:
//  1. X-Bundler-Client-Id header: See UseClientIDHeader
//  2. X-Forwarded-By header: The first IP address in the array which is assumed to be the client
//  3. Request.RemoteAddr: The remote IP address
func (r *Relayer) FilterByClient(c *gin.Context) {
	isBanned := false
	err := r.db.View(func(txn *badger.Txn) error {
		opsSeen, opsIncluded, err := getOpsCountByClientID(txn, r.getClientID(c))
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
		r.errorHandler(err)
		c.Status(http.StatusInternalServerError)
		c.Abort()
	}

	if isBanned {
		c.Status(http.StatusForbidden)
		c.Abort()
	}
}

// LogClientForSendUserOperation is a custom Gin middleware used to map a userOp requestID to a client
// identifier (e.g. IP address).
func (r *Relayer) LogClientForSendUserOperation(c *gin.Context) {
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
		r.errorHandler(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	rid := op.GetRequestID(common.HexToAddress(ep), r.chainID).String()
	cid := r.getClientID(c)
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
		r.errorHandler(err)
		c.Status(http.StatusInternalServerError)
		return
	}
}

// BatchHandler accepts a batch and sends it as a regular EOA transaction.
func (r *Relayer) BatchHandler(ctx *modules.BatchHandlerCtx) error {
	err := r.db.Update(func(txn *badger.Txn) error {
		if len(ctx.PendingRemoval) > 0 {
			rids := getRequestIDsFromOps(ctx.EntryPoint, ctx.ChainID, ctx.PendingRemoval...)
			if err := removeRequestIDEntries(txn, rids...); err != nil {
				return err
			}
		}

		if len(ctx.Batch) > 0 {
			// TODO: Submit batch and handle reverts.

			rids := getRequestIDsFromOps(ctx.EntryPoint, ctx.ChainID, ctx.Batch...)
			if err := incrementOpsIncludedByRequestIDs(txn, rids...); err != nil {
				return err
			}
			if err := removeRequestIDEntries(txn, rids...); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
