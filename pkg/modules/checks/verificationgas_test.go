package checks

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
)

// TestOpVGlessThanMaxVG calls checks.ValidateVerificationGas where verificationGas < MAX_VERIFICATION_GAS.
// Expects nil.
func TestOpVGlessThanMaxVG(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	mvg := big.NewInt(0).Add(op.VerificationGasLimit, common.Big1)

	if err := ValidateVerificationGas(op, mvg); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestOpVGEqualMaxVG calls checks.ValidateVerificationGas where verificationGas == MAX_VERIFICATION_GAS.
// Expects nil.
func TestOpVGEqualMaxVG(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	mvg := big.NewInt(0).Add(op.VerificationGasLimit, common.Big0)

	if err := ValidateVerificationGas(op, mvg); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestOpVGMoreThanMaxVG calls checks.ValidateVerificationGas where verificationGas > MAX_VERIFICATION_GAS.
// Expects error.
func TestOpVGMoreThanMaxVG(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	mvg := big.NewInt(0).Sub(op.VerificationGasLimit, common.Big1)

	if err := ValidateVerificationGas(op, mvg); err == nil {
		t.Fatal("got nil, want err")
	}
}

// TestOpPVGMoreThanOH calls checks.ValidateVerificationGas where the preVerificationGas > overhead gas.
// Expect nil.
func TestOpPVGMoreThanOH(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	pvg := gas.NewDefaultOverhead().CalcPreVerificationGas(op)
	op.PreVerificationGas = big.NewInt(0).Add(pvg, common.Big1)

	if err := ValidateVerificationGas(op, op.VerificationGasLimit); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestOpPVGEqualOH calls checks.ValidateVerificationGas where the preVerificationGas == overhead gas. Expect
// nil.
func TestOpPVGEqualOH(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	pvg := gas.NewDefaultOverhead().CalcPreVerificationGas(op)
	op.PreVerificationGas = big.NewInt(0).Add(pvg, common.Big0)

	if err := ValidateVerificationGas(op, op.VerificationGasLimit); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestOpPVGLessThanOH calls checks.ValidateVerificationGas where the preVerificationGas < overhead gas.
// Expect error.
func TestOpPVGLessThanOH(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	pvg := gas.NewDefaultOverhead().CalcPreVerificationGas(op)
	op.PreVerificationGas = big.NewInt(0).Sub(pvg, common.Big1)

	if err := ValidateVerificationGas(op, op.VerificationGasLimit); err == nil {
		t.Fatal("got nil, want err")
	}
}
