package checks

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
)

// TestOpCGLessThanOH calls checks.ValidateCallGasLimit where callGas < overhead. Expect error.
func TestOpCGLessThanOH(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	cg := gas.NewDefaultOverhead().NonZeroValueCall()
	op.CallGasLimit = big.NewInt(0).Sub(cg, common.Big1)
	err := ValidateCallGasLimit(op)

	if err == nil {
		t.Fatalf("got nil, want err")
	}
}

// TestOpCGEqualOH calls checks.ValidateCallGasLimit where callGas == overhead. Expect nil.
func TestOpCGEqualOH(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	cg := gas.NewDefaultOverhead().NonZeroValueCall()
	op.CallGasLimit = big.NewInt(0).Add(cg, common.Big0)
	err := ValidateCallGasLimit(op)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestOpCGMoreThanOH calls checks.ValidateCallGasLimit where callGas > overhead. Expect nil.
func TestOpCGMoreThanOH(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	cg := gas.NewDefaultOverhead().NonZeroValueCall()
	op.CallGasLimit = big.NewInt(0).Add(cg, common.Big1)
	err := ValidateCallGasLimit(op)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}
