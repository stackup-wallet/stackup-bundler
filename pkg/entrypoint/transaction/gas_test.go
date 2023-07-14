package transaction

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// TestSuggestMeanGasTipCapForNormalLoad simulates a scenario of normal network load. In this case the average
// tip from userOps is assumed to not be above the suggested tip. In which case the suggested tip is the
// optimal choice.
func TestSuggestMeanGasTipCapForNormalLoad(t *testing.T) {
	op1 := testutils.MockValidInitUserOp()
	op1.MaxPriorityFeePerGas = big.NewInt(1)
	op2 := testutils.MockValidInitUserOp()
	op2.MaxPriorityFeePerGas = big.NewInt(1)
	batch := []*userop.UserOperation{op1, op2}

	expected := big.NewInt(2)
	if tip := SuggestMeanGasTipCap(expected, batch); tip.Cmp(expected) != 0 {
		t.Fatalf("got %d, want %d", tip.Int64(), expected.Int64())
	}
}

// TestSuggestMeanGasTipCapForHighLoad simulates a scenario of high network load. In this case the average tip
// from userOps is assumed to be above the suggested tip (i.e. userOps want to be included quickly). In which
// case the average tip from userOps is the optimal choice.
func TestSuggestMeanGasTipCapForHighLoad(t *testing.T) {
	op1 := testutils.MockValidInitUserOp()
	op1.MaxPriorityFeePerGas = big.NewInt(5)
	op2 := testutils.MockValidInitUserOp()
	op2.MaxPriorityFeePerGas = big.NewInt(10)
	batch := []*userop.UserOperation{op1, op2}

	expected := big.NewInt(7)
	if tip := SuggestMeanGasTipCap(big.NewInt(2), batch); tip.Cmp(expected) != 0 {
		t.Fatalf("got %d, want %d", tip.Int64(), expected.Int64())
	}
}

// TestSuggestMeanGasFeeCapNormalLoad simulates a scenario of normal network load. In this case the average
// gas fee cap from userOps is assumed to not be above the recommended max fee. In which case the recommended
// max fee is the optimal choice.
func TestSuggestMeanGasFeeCapNormalLoad(t *testing.T) {
	op1 := testutils.MockValidInitUserOp()
	op1.MaxFeePerGas = big.NewInt(1)
	op2 := testutils.MockValidInitUserOp()
	op2.MaxFeePerGas = big.NewInt(1)
	batch := []*userop.UserOperation{op1, op2}

	bf := big.NewInt(1)
	tip := big.NewInt(0)
	expected := big.NewInt(0).Mul(bf, common.Big2)
	mf := SuggestMeanGasFeeCap(bf, tip, batch)
	if mf.Cmp(expected) != 0 {
		t.Fatalf("got %d, want %d", mf.Int64(), expected.Int64())
	}
}

// TestSuggestMeanGasFeeCapHighLoad simulates a scenario of high network load. In this case the average gas
// fee cap from userOps is assumed to be above the recommended max fee. In which case the average gas fee cap
// from userOps is the optimal choice.
func TestSuggestMeanGasFeeCapHighLoad(t *testing.T) {
	op1 := testutils.MockValidInitUserOp()
	op1.MaxFeePerGas = big.NewInt(5)
	op2 := testutils.MockValidInitUserOp()
	op2.MaxFeePerGas = big.NewInt(10)
	batch := []*userop.UserOperation{op1, op2}

	bf := big.NewInt(1)
	tip := big.NewInt(0)
	expected := big.NewInt(7)
	mf := SuggestMeanGasFeeCap(bf, tip, batch)
	if mf.Cmp(expected) != 0 {
		t.Fatalf("got %d, want %d", mf.Int64(), expected.Int64())
	}
}

// TestSuggestMeanGasPriceForNormalLoad simulates a scenario of normal network load. In this case the average
// gas fee cap from userOps is assumed to not be above the given gas price. In which case the given gas price
// is the optimal choice.
func TestSuggestMeanGasPriceForNormalLoad(t *testing.T) {
	op1 := testutils.MockValidInitUserOp()
	op1.MaxFeePerGas = big.NewInt(1)
	op2 := testutils.MockValidInitUserOp()
	op2.MaxFeePerGas = big.NewInt(1)

	expected := big.NewInt(2)
	batch := []*userop.UserOperation{op1, op2}
	gp := SuggestMeanGasPrice(expected, batch)
	if gp.Cmp(expected) != 0 {
		t.Fatalf("got %d, want %d", gp.Int64(), expected.Int64())
	}
}

// TestSuggestMeanGasPriceForHighLoad simulates a scenario of high network load. In this case the average gas
// fee cap from userOps is assumed to be above the given gas price. In which case the average gas fee cap from
// userOps is the optimal choice.
func TestSuggestMeanGasPriceForHighLoad(t *testing.T) {
	op1 := testutils.MockValidInitUserOp()
	op1.MaxFeePerGas = big.NewInt(5)
	op2 := testutils.MockValidInitUserOp()
	op2.MaxFeePerGas = big.NewInt(10)

	expected := big.NewInt(7)
	batch := []*userop.UserOperation{op1, op2}
	gp := SuggestMeanGasPrice(big.NewInt(2), batch)
	if gp.Cmp(expected) != 0 {
		t.Fatalf("got %d, want %d", gp.Int64(), expected.Int64())
	}
}
