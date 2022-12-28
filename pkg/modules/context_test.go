package modules

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// TestAddDepositInfoToCtx verifies that stake info can be added to a context and later retrieved.
func TestAddDepositInfoToCtx(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{}
	ctx := NewUserOpHandlerContext(op, penOps, testutils.ValidAddress, testutils.ChainID)

	entity := op.GetFactory()
	dep := testutils.StakedDepositInfo
	ctx.AddDepositInfo(entity, dep)

	if ctx.GetDepositInfo(entity) != dep {
		t.Fatal("Retrieved deposit info does not equal the original")
	}
}

// TestGetNilDepositInfoFromCtx calls (c *UserOpHandlerCtx).GetDepositInfo on an address that has not been
// set. Expects nil.
func TestGetNilDepositInfoFromCtx(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{}
	ctx := NewUserOpHandlerContext(op, penOps, testutils.ValidAddress, testutils.ChainID)

	if dep := ctx.GetDepositInfo(op.GetFactory()); dep != nil {
		t.Fatalf("got %+v, want nil", dep)
	}
}

// TestGetPendingOps calls (c *UserOpHandlerCtx).GetPendingOps and verifies that it returns the same array of
// UserOperations the context was initialized with.
func TestGetPendingOps(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	penOp1 := testutils.MockValidInitUserOp()
	penOp2 := testutils.MockValidInitUserOp()
	penOp2.Nonce = big.NewInt(0).Add(penOp1.Nonce, common.Big1)
	penOp3 := testutils.MockValidInitUserOp()
	penOp3.Nonce = big.NewInt(0).Add(penOp2.Nonce, common.Big1)
	initPenOps := []*userop.UserOperation{penOp1, penOp2, penOp3}
	ctx := NewUserOpHandlerContext(op, initPenOps, testutils.ValidAddress, testutils.ChainID)

	penOps := ctx.GetPendingOps()
	if len(penOps) != len(initPenOps) {
		t.Fatalf("got length %d, want %d", len(penOps), len(initPenOps))
	}

	for i, penOp := range penOps {
		if !testutils.IsOpsEqual(penOp, initPenOps[i]) {
			t.Fatalf("ops not equal: %s", testutils.GetOpsDiff(penOp, initPenOps[i]))
		}
	}
}
