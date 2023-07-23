package checks

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

// TestOpMAGLessThanMax calls checks.ValidateGasAvailable where op.GetMaxAvailableGas < MaxBatchGasLimit.
// Expect nil.
func TestOpMAGLessThanMax(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	max := big.NewInt(0).Add(op.GetMaxGasAvailable(), common.Big1)
	err := ValidateGasAvailable(op, max)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestOpMAGEqualToMax calls checks.ValidateGasAvailable where op.GetMaxAvailableGas == MaxBatchGasLimit.
// Expect nil.
func TestOpMAGEqualToMax(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	err := ValidateGasAvailable(op, op.GetMaxGasAvailable())

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestOpMAGMoreThanMax calls checks.ValidateGasAvailable where op.GetMaxAvailableGas > MaxBatchGasLimit.
// Expect error.
func TestOpMAGMoreThanMax(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	max := big.NewInt(0).Sub(op.GetMaxGasAvailable(), common.Big1)
	err := ValidateGasAvailable(op, max)

	if err == nil {
		t.Fatalf("got nil, want err")
	}
}
