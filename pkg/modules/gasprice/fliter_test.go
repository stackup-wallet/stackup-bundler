package gasprice_test

import (
	"math/big"
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/gasprice"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// TestFilterUnderpricedDynamic verifies that FilterUnderpriced will remove all UserOperations from a batch
// where the effective gas price is less than the expected bundler transaction's gas price.
func TestFilterUnderpricedDynamic(t *testing.T) {
	bf := big.NewInt(4)
	tip := big.NewInt(1)

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

	ctx := modules.NewBatchHandlerContext(
		[]*userop.UserOperation{op1, op2, op3},
		testutils.ValidAddress1,
		testutils.ChainID,
		bf,
		tip,
		big.NewInt(10),
	)
	if err := gasprice.FilterUnderpriced()(ctx); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if len(ctx.Batch) != 2 {
		t.Fatalf("got length %d, want 2", len(ctx.Batch))
	} else if !testutils.IsOpsEqual(ctx.Batch[0], op2) {
		t.Fatal("incorrect order: first op out of place")
	} else if !testutils.IsOpsEqual(ctx.Batch[1], op3) {
		t.Fatal("incorrect order: second op out of place")
	}
}

// TestFilterUnderpricedGasPrice verifies that FilterUnderpriced will remove all UserOperations from a batch
// where the MaxFeePerGas is less than the context GasPrice.
func TestFilterUnderpricedGasPrice(t *testing.T) {
	op1 := testutils.MockValidInitUserOp()
	op1.MaxFeePerGas = big.NewInt(4)
	op1.MaxPriorityFeePerGas = big.NewInt(4)

	op2 := testutils.MockValidInitUserOp()
	op2.Sender = testutils.ValidAddress2
	op2.MaxFeePerGas = big.NewInt(5)
	op2.MaxPriorityFeePerGas = big.NewInt(5)

	op3 := testutils.MockValidInitUserOp()
	op3.Sender = testutils.ValidAddress3
	op3.MaxFeePerGas = big.NewInt(6)
	op3.MaxPriorityFeePerGas = big.NewInt(6)

	ctx := modules.NewBatchHandlerContext(
		[]*userop.UserOperation{op1, op2, op3},
		testutils.ValidAddress1,
		testutils.ChainID,
		nil,
		nil,
		big.NewInt(5),
	)
	if err := gasprice.FilterUnderpriced()(ctx); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if len(ctx.Batch) != 2 {
		t.Fatalf("got length %d, want 2", len(ctx.Batch))
	} else if !testutils.IsOpsEqual(ctx.Batch[0], op2) {
		t.Fatal("incorrect order: first op out of place")
	} else if !testutils.IsOpsEqual(ctx.Batch[1], op3) {
		t.Fatal("incorrect order: second op out of place")
	}
}
