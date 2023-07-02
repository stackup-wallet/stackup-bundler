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
	ep := testutils.ValidAddress1
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
	ep := testutils.ValidAddress1
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
	ep := testutils.ValidAddress1
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

// TestDumpFromMempool verifies that bundles are being built with UserOperations in the mempool. Ordering is
// FIFO and more specific sorting and filtering is left up to downstream modules to implement.
func TestDumpFromMempool(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem, _ := New(db)
	ep := testutils.ValidAddress1

	op1 := testutils.MockValidInitUserOp()
	op1.MaxFeePerGas = big.NewInt(4)
	op1.MaxPriorityFeePerGas = big.NewInt(3)

	op2 := testutils.MockValidInitUserOp()
	op2.Sender = testutils.ValidAddress2
	op2.MaxFeePerGas = big.NewInt(5)
	op2.MaxPriorityFeePerGas = big.NewInt(2)

	op3 := testutils.MockValidInitUserOp()
	op3.Sender = testutils.ValidAddress3
	op3.MaxFeePerGas = big.NewInt(6)
	op3.MaxPriorityFeePerGas = big.NewInt(1)

	if err := mem.AddOp(ep, op1); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if err := mem.AddOp(ep, op2); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	if err := mem.AddOp(ep, op3); err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	if memOps, err := mem.Dump(ep); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if len(memOps) != 3 {
		t.Fatalf("got length %d, want 3", len(memOps))
	} else if !testutils.IsOpsEqual(memOps[0], op1) {
		t.Fatal("incorrect order: first op out of place")
	} else if !testutils.IsOpsEqual(memOps[1], op2) {
		t.Fatal("incorrect order: second op out of place")
	} else if !testutils.IsOpsEqual(memOps[2], op3) {
		t.Fatal("incorrect order: third op out of place")
	}
}

// TestNewMempoolLoadsFromDisk verifies that a new Mempool instance is built from ops saved in the DB without
// including ops previously removed.
func TestNewMempoolLoadsFromDisk(t *testing.T) {
	db := testutils.DBMock()
	defer db.Close()
	mem1, _ := New(db)
	ep := testutils.ValidAddress1
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
