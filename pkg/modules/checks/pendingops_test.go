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
	err := ValidatePendingOps(
		op,
		penOps,
		testutils.MaxOpsForUnstakedSender,
		testutils.MockGetNotStakeZeroDeposit,
	)

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
	err := ValidatePendingOps(
		op,
		penOps,
		testutils.MaxOpsForUnstakedSender,
		testutils.MockGetNotStakeZeroDeposit,
	)

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
	err := ValidatePendingOps(
		op,
		penOps,
		testutils.MaxOpsForUnstakedSender,
		testutils.MockGetStakeZeroDeposit,
	)

	if err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}
