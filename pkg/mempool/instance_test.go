package mempool

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

// TestAddOpToMempool verifies that a UserOperation can be added to the mempool and later retrieved without
// any changes.
func TestAddOpToMempool(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := New(db)
	ep := testutils.ValidAddress
	op := testutils.MockValidInitUserOp()

	if err := mem.AddOp(ep, op); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	memOps, err := mem.GetOps(ep, op.Sender)
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if len(memOps) != 1 {
		t.Fatalf("got length %d, want 1", len(memOps))
	}

	if !testutils.IsOpsEqual(op, memOps[0]) {
		t.Fatalf("ops not equal: %s", testutils.GetOpsDiff(op, memOps[0]))
	}
}

// TestReplaceOpInMempool verifies that a UserOperation with same Sender and Nonce can replace another
// UserOperation already in the mempool.
func TestReplaceOpInMempool(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := New(db)
	ep := testutils.ValidAddress
	op1 := testutils.MockValidInitUserOp()
	op2 := testutils.MockValidInitUserOp()
	op2.MaxPriorityFeePerGas = big.NewInt(0).Add(op1.MaxPriorityFeePerGas, common.Big1)

	if err := mem.AddOp(ep, op1); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if err := mem.AddOp(ep, op2); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	memOps, err := mem.GetOps(ep, op2.Sender)
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if len(memOps) != 1 {
		t.Fatalf("got length %d, want 1", len(memOps))
	}

	if !testutils.IsOpsEqual(op2, memOps[0]) {
		t.Fatalf("ops not equal: %s", testutils.GetOpsDiff(op2, memOps[0]))
	}
}

// TestRemoveOpsFromMempool verifies that a UserOperation can be added to the mempool and later removed.
func TestRemoveOpsFromMempool(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := New(db)
	ep := testutils.ValidAddress
	op := testutils.MockValidInitUserOp()

	if err := mem.AddOp(ep, op); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	if err := mem.RemoveOps(ep, op); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	memOps, err := mem.GetOps(ep, op.Sender)
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if len(memOps) != 0 {
		t.Fatalf("got length %d, want 0", len(memOps))
	}
}

// TestBundleOpsFromMempool verifies that bundles are being built with the correct ordering of highest
// MaxPriorityFeePerGas first.
func TestBundleOpsFromMempool(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := New(db)
	ep := testutils.ValidAddress
	op1 := testutils.MockValidInitUserOp()
	op2 := testutils.MockValidInitUserOp()
	op2.Nonce = big.NewInt(0).Add(op1.Nonce, common.Big1)
	op2.MaxPriorityFeePerGas = big.NewInt(0).Add(op1.MaxPriorityFeePerGas, common.Big1)

	if err := mem.AddOp(ep, op1); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if err := mem.AddOp(ep, op2); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	memOps, err := mem.BundleOps(ep)
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if len(memOps) != 2 {
		t.Fatalf("got length %d, want 2", len(memOps))
	}

	if !testutils.IsOpsEqual(op2, memOps[0]) {
		t.Fatalf("incorrect order, expect ops with highest MaxPriorityFeePerGas first")
	}
}

// TestNewMempoolLoadsFromDisk verifies that a new Mempool instance is built from ops saved in the DB without
// including ops previously removed.
func TestNewMempoolLoadsFromDisk(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem1, _ := New(db)
	ep := testutils.ValidAddress
	op1 := testutils.MockValidInitUserOp()
	op2 := testutils.MockValidInitUserOp()
	op2.Nonce = big.NewInt(0).Add(op1.Nonce, common.Big1)
	op2.MaxPriorityFeePerGas = big.NewInt(0).Add(op1.MaxPriorityFeePerGas, common.Big1)

	if err := mem1.AddOp(ep, op1); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if err := mem1.AddOp(ep, op2); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if err := mem1.RemoveOps(ep, op1); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	mem2, _ := New(db)
	memOps, err := mem2.GetOps(ep, op2.Sender)
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if len(memOps) != 1 {
		t.Fatalf("got length %d, want 1", len(memOps))
	}

	if !testutils.IsOpsEqual(op2, memOps[0]) {
		t.Fatalf("ops not equal: %s", testutils.GetOpsDiff(op2, memOps[0]))
	}
}
