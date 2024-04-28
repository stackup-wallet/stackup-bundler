package checks

import (
	"errors"
	"math/big"
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func TestNoPendingOps(t *testing.T) {
	penOps := []*userop.UserOperationV06{}
	op := testutils.MockValidInitV06UserOp()
	err := ValidatePendingOps(
		op,
		penOps,
	)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}

func TestPendingOpsWithNewOp(t *testing.T) {
	penOp := testutils.MockValidInitV06UserOp()
	penOps := []*userop.UserOperationV06{penOp}
	op := testutils.MockValidInitV06UserOp()
	op.Nonce = big.NewInt(1)
	err := ValidatePendingOps(
		op,
		penOps,
	)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}

func TestPendingOpsWithNoGasFeeReplacement(t *testing.T) {
	penOp := testutils.MockValidInitV06UserOp()
	penOps := []*userop.UserOperationV06{penOp}
	op := testutils.MockValidInitV06UserOp()
	err := ValidatePendingOps(
		op,
		penOps,
	)

	if !errors.Is(err, ErrReplacementOpUnderpriced) {
		t.Fatalf("got %v, want ErrReplacementOpUnderpriced", err)
	}
}

func TestPendingOpsWithOnlyMaxFeeReplacement(t *testing.T) {
	penOp := testutils.MockValidInitV06UserOp()
	penOps := []*userop.UserOperationV06{penOp}
	op := testutils.MockValidInitV06UserOp()
	op.MaxFeePerGas, _ = calcNewThresholds(op.MaxFeePerGas, op.MaxPriorityFeePerGas)
	err := ValidatePendingOps(
		op,
		penOps,
	)

	if !errors.Is(err, ErrReplacementOpUnderpriced) {
		t.Fatalf("got %v, want ErrReplacementOpUnderpriced", err)
	}
}

func TestPendingOpsWithOnlyMaxPriorityFeeReplacement(t *testing.T) {
	penOp := testutils.MockValidInitV06UserOp()
	penOps := []*userop.UserOperationV06{penOp}
	op := testutils.MockValidInitV06UserOp()
	_, op.MaxPriorityFeePerGas = calcNewThresholds(op.MaxFeePerGas, op.MaxPriorityFeePerGas)
	err := ValidatePendingOps(
		op,
		penOps,
	)

	if !errors.Is(err, ErrReplacementOpUnderpriced) {
		t.Fatalf("got %v, want ErrReplacementOpUnderpriced", err)
	}
}

func TestPendingOpsWithOkGasFeeReplacement(t *testing.T) {
	penOp := testutils.MockValidInitV06UserOp()
	penOps := []*userop.UserOperationV06{penOp}
	op := testutils.MockValidInitV06UserOp()
	op.MaxFeePerGas, op.MaxPriorityFeePerGas = calcNewThresholds(op.MaxFeePerGas, op.MaxPriorityFeePerGas)
	err := ValidatePendingOps(
		op,
		penOps,
	)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}
