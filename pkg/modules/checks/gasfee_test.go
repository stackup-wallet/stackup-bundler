package checks

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

// TestMPFLessThanGasTip calls checks.ValidateFeePerGas with a maxPriorityFeePerGas < gas tip. Expect error.
func TestMPFLessThanGasTip(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	gt := testutils.GetMockGasTipFunc(big.NewInt(0).Add(op.MaxPriorityFeePerGas, common.Big1))
	err := ValidateFeePerGas(op, gt)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}

// TestMPFEqualGasTip calls checks.ValidateFeePerGas with a maxPriorityFeePerGas == gas tip. Expect nil.
func TestMPFEqualGasTip(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	gt := testutils.GetMockGasTipFunc(big.NewInt(0).Add(op.MaxPriorityFeePerGas, common.Big0))
	err := ValidateFeePerGas(op, gt)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestMPFMoreThanGasTip calls checks.ValidateFeePerGas with a maxPriorityFeePerGas > gas tip. Expect nil.
func TestMPFMoreThanGasTip(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	gt := testutils.GetMockGasTipFunc(big.NewInt(0).Sub(op.MaxPriorityFeePerGas, common.Big1))
	err := ValidateFeePerGas(op, gt)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestMFLessThanMPF calls checks.ValidateFeePerGas with a MaxFeePerGas < maxPriorityFeePerGas. Expect error.
func TestMFLessThanMPF(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	gt := testutils.GetMockGasTipFunc(big.NewInt(0).Add(op.MaxPriorityFeePerGas, common.Big0))
	op.MaxFeePerGas = big.NewInt(0).Sub(op.MaxPriorityFeePerGas, common.Big1)
	err := ValidateFeePerGas(op, gt)
	fmt.Println(err)

	if err == nil {
		t.Fatal("got nil, want err")
	}
}

// TestMFEqualMPF calls checks.ValidateFeePerGas with a MaxFeePerGas == maxPriorityFeePerGas. Expect nil.
func TestMFEqualMPF(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	gt := testutils.GetMockGasTipFunc(big.NewInt(0).Add(op.MaxPriorityFeePerGas, common.Big0))
	op.MaxFeePerGas = big.NewInt(0).Add(op.MaxPriorityFeePerGas, common.Big0)
	err := ValidateFeePerGas(op, gt)
	fmt.Println(err)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

// TestMFMoreThanMPF calls checks.ValidateFeePerGas with a MaxFeePerGas > maxPriorityFeePerGas. Expect nil.
func TestMFMoreThanMPF(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	gt := testutils.GetMockGasTipFunc(big.NewInt(0).Add(op.MaxPriorityFeePerGas, common.Big0))
	op.MaxFeePerGas = big.NewInt(0).Add(op.MaxPriorityFeePerGas, common.Big1)
	err := ValidateFeePerGas(op, gt)
	fmt.Println(err)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}
