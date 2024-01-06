package checks

import (
	"errors"
	"math/big"
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func TestNoPendingOps(t *testing.T) {
	penOps := []*userop.UserOperation{}
	op := testutils.MockValidInitUserOp()
	err := ValidatePendingOps(
		op,
		penOps,
	)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}

func TestPendingOpsWithNewOp(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	op.Nonce = big.NewInt(1)
	err := ValidatePendingOps(
		op,
		penOps,
	)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}

func TestPendingOpsWithFailGasFeeReplacement(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	err := ValidatePendingOps(
		op,
		penOps,
	)

	if !errors.Is(err, ErrReplacementOpUnderpriced) {
		t.Fatalf("got %v, want ErrReplacementOpUnderpriced", err)
	}
}

func TestPendingOpsWithFailMaxFeeReplacement(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	_, op.MaxPriorityFeePerGas = calcNewThresholds(op.MaxFeePerGas, op.MaxPriorityFeePerGas)
	err := ValidatePendingOps(
		op,
		penOps,
	)

	if !errors.Is(err, ErrReplacementOpUnderpriced) {
		t.Fatalf("got %v, want ErrReplacementOpUnderpriced", err)
	}
}

func TestPendingOpsWithFailMaxPriorityFeeReplacement(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	op.MaxFeePerGas, _ = calcNewThresholds(op.MaxFeePerGas, op.MaxPriorityFeePerGas)
	err := ValidatePendingOps(
		op,
		penOps,
	)

	if !errors.Is(err, ErrReplacementOpUnderpriced) {
		t.Fatalf("got %v, want ErrReplacementOpUnderpriced", err)
	}
}

func TestPendingOpsWithOkGasFeeReplacement(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	op.MaxFeePerGas, op.MaxPriorityFeePerGas = calcNewThresholds(op.MaxFeePerGas, op.MaxPriorityFeePerGas)
	err := ValidatePendingOps(
		op,
		penOps,
	)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}
