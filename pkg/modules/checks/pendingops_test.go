package checks

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// TestNoPendingOps calls checks.ValidatePendingOps with no pending UserOperations. Expect nil.
func TestNoPendingOps(t *testing.T) {
	penOps := []*userop.UserOperation{}
	op := testutils.MockValidInitUserOp()
	err := ValidatePendingOps(op, penOps, testutils.MockGetNotStakeZeroDeposit)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}

// TestPendingOpsNotStaked calls checks.ValidatePendingOps with pending UserOperations but sender is not
// staked. Expect error.
func TestPendingOpsNotStaked(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	op.Nonce = big.NewInt(0).Add(penOp.Nonce, common.Big1)
	err := ValidatePendingOps(op, penOps, testutils.MockGetNotStakeZeroDeposit)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}

// TestPendingOpsStaked calls checks.ValidatePendingOps with pending UserOperations but sender is staked.
// Expect nil.
func TestPendingOpsStaked(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	op.Nonce = big.NewInt(0).Add(penOp.Nonce, common.Big1)
	err := ValidatePendingOps(op, penOps, testutils.MockGetStakeZeroDeposit)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}

// TestReplaceOp calls checks.ValidatePendingOps with a valid UserOperation that replaces a pending
// UserOperation. Expect nil.
func TestReplaceOp(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	op.MaxPriorityFeePerGas = big.NewInt(0).Add(penOp.MaxPriorityFeePerGas, common.Big1)
	op.MaxFeePerGas = big.NewInt(0).Add(penOp.MaxFeePerGas, common.Big1)
	err := ValidatePendingOps(op, penOps, testutils.MockGetNotStakeZeroDeposit)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}

// TestReplaceOpLowerMPF calls checks.ValidatePendingOps with a UserOperation that replaces a pending
// UserOperation but has a lower MaxPriorityFeePerGas. Expect error.
func TestReplaceOpLowerMPF(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	op.MaxPriorityFeePerGas = big.NewInt(0).Sub(penOp.MaxPriorityFeePerGas, common.Big1)
	err := ValidatePendingOps(op, penOps, testutils.MockGetNotStakeZeroDeposit)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}

// TestReplaceOpEqualMPF calls checks.ValidatePendingOps with a UserOperation that replaces a pending
// UserOperation but has an equal MaxPriorityFeePerGas. Expect error.
func TestReplaceOpEqualMPF(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	op.MaxPriorityFeePerGas = big.NewInt(0).Add(penOp.MaxPriorityFeePerGas, common.Big0)
	err := ValidatePendingOps(op, penOps, testutils.MockGetNotStakeZeroDeposit)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}

// TestReplaceOpNotEqualIncMF calls checks.ValidatePendingOps with a UserOperation that replaces a pending
// UserOperation but does not have an equally increasing MaxFeePerGas. Expect error.
func TestReplaceOpNotEqualIncMF(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	op.MaxPriorityFeePerGas = big.NewInt(0).Add(penOp.MaxPriorityFeePerGas, common.Big2)
	op.MaxFeePerGas = big.NewInt(0).Add(penOp.MaxFeePerGas, common.Big1)
	err := ValidatePendingOps(op, penOps, testutils.MockGetNotStakeZeroDeposit)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}

// TestReplaceOpSameMF calls checks.ValidatePendingOps with a UserOperation that replaces a pending
// UserOperation but does not increase MaxFeePerGas. Expect error.
func TestReplaceOpSameMF(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	op.MaxPriorityFeePerGas = big.NewInt(0).Add(penOp.MaxPriorityFeePerGas, common.Big1)
	op.MaxFeePerGas = big.NewInt(0).Add(penOp.MaxFeePerGas, common.Big0)
	err := ValidatePendingOps(op, penOps, testutils.MockGetNotStakeZeroDeposit)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}

// TestReplaceOpDecMF calls checks.ValidatePendingOps with a UserOperation that replaces a pending
// UserOperation but has a decreasing MaxFeePerGas. Expect error.
func TestReplaceOpDecMF(t *testing.T) {
	penOp := testutils.MockValidInitUserOp()
	penOps := []*userop.UserOperation{penOp}
	op := testutils.MockValidInitUserOp()
	op.MaxPriorityFeePerGas = big.NewInt(0).Add(penOp.MaxPriorityFeePerGas, common.Big1)
	op.MaxFeePerGas = big.NewInt(0).Sub(penOp.MaxFeePerGas, common.Big1)
	err := ValidatePendingOps(op, penOps, testutils.MockGetNotStakeZeroDeposit)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}
