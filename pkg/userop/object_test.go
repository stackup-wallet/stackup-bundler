package userop_test

import (
	"math/big"
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

// TestUserOperationGetDynamicGasPrice verifies that (*UserOperation).GetDynamicGasPrice returns the correct
// effective gas price given a base fee.
func TestUserOperationGetDynamicGasPrice(t *testing.T) {
	bf := big.NewInt(3)
	op := testutils.MockValidInitUserOp()

	// basefee + MPF > MF
	want := big.NewInt(4)
	op.MaxFeePerGas = big.NewInt(4)
	op.MaxPriorityFeePerGas = big.NewInt(3)
	if op.GetDynamicGasPrice(bf).Cmp(want) != 0 {
		t.Fatalf("got %d, want %d", op.GetDynamicGasPrice(bf).Int64(), want.Int64())
	}

	// basefee + MPF == MF
	want = big.NewInt(5)
	op.MaxFeePerGas = big.NewInt(5)
	op.MaxPriorityFeePerGas = big.NewInt(2)
	if op.GetDynamicGasPrice(bf).Cmp(want) != 0 {
		t.Fatalf("got %d, want %d", op.GetDynamicGasPrice(bf).Int64(), want.Int64())
	}

	// basefee + MPF < MF
	want = big.NewInt(4)
	op.MaxFeePerGas = big.NewInt(6)
	op.MaxPriorityFeePerGas = big.NewInt(1)
	if op.GetDynamicGasPrice(bf).Cmp(want) != 0 {
		t.Fatalf("got %d, want %d", op.GetDynamicGasPrice(bf).Int64(), want.Int64())
	}
}

// TestUserOperationGetGasPriceNilBF verifies that (*UserOperation).GetDynamicGasPrice returns the correct
// value when basefee is nil.
func TestUserOperationGetGasPriceNilBF(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	op.MaxFeePerGas = big.NewInt(4)
	op.MaxPriorityFeePerGas = big.NewInt(3)
	if op.GetDynamicGasPrice(nil).Cmp(op.MaxPriorityFeePerGas) != 0 {
		t.Fatalf("got %d, want %d", op.GetDynamicGasPrice(nil).Int64(), op.MaxPriorityFeePerGas)
	}
}
