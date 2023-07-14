package checks

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

// TestMFLessThanBF calls checks.ValidateFeePerGas with a MaxFeePerGas < base fee. Expect error.
func TestMFLessThanBF(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	gbf := testutils.GetMockBaseFeeFunc(common.Big2)
	op.MaxFeePerGas = common.Big1
	op.MaxPriorityFeePerGas = common.Big0
	err := ValidateFeePerGas(op, gbf)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}

// TestMFEqualBF calls checks.ValidateFeePerGas with a MaxFeePerGas == base fee. Expect nil.
func TestMFEqualBF(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	gbf := testutils.GetMockBaseFeeFunc(common.Big1)
	op.MaxFeePerGas = common.Big1
	op.MaxPriorityFeePerGas = common.Big0
	err := ValidateFeePerGas(op, gbf)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestMFMoreThanBF calls checks.ValidateFeePerGas with a MaxFeePerGas > base fee. Expect nil.
func TestMFMoreThanBF(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	gbf := testutils.GetMockBaseFeeFunc(common.Big1)
	op.MaxFeePerGas = common.Big2
	op.MaxPriorityFeePerGas = common.Big0
	err := ValidateFeePerGas(op, gbf)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestMPFMoreThanMF calls checks.ValidateFeePerGas with a MaxPriorityFeePerGas > MaxFeePerGas. Expect error.
func TestMPFMoreThanMF(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	gbf := testutils.GetMockBaseFeeFunc(common.Big1)
	op.MaxFeePerGas = common.Big2
	op.MaxPriorityFeePerGas = common.Big3
	err := ValidateFeePerGas(op, gbf)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}

// TestMPFEqualMF calls checks.ValidateFeePerGas with a MaxPriorityFeePerGas == MaxFeePerGas. Expect nil.
func TestMPFEqualMF(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	gbf := testutils.GetMockBaseFeeFunc(common.Big1)
	op.MaxFeePerGas = common.Big2
	op.MaxPriorityFeePerGas = common.Big2
	err := ValidateFeePerGas(op, gbf)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestMPFLessThanMF calls checks.ValidateFeePerGas with a MaxPriorityFeePerGas < MaxFeePerGas. Expect nil.
func TestMPFLessThanMF(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	gbf := testutils.GetMockBaseFeeFunc(common.Big1)
	op.MaxFeePerGas = common.Big2
	op.MaxPriorityFeePerGas = common.Big1
	err := ValidateFeePerGas(op, gbf)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}
