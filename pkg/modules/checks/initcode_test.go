package checks

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
)

// TestInitCodeDNE calls checks.ValidateInitCode where initCode does not exist. Expect nil.
func TestInitCodeDNE(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	op.InitCode = []byte{}
	err := ValidateInitCode(op, func(f common.Address) (*entrypoint.IStakeManagerDepositInfo, error) {
		return nil, nil
	})

	if err != nil {
		t.Fatalf(`got err %v, want nil`, err)
	}
}

// TestInitCodeContainsAddress calls checks.ValidateInitCode where initCode exist without a valid address.
// Expect error.
func TestInitCodeContainsAddress(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	op.InitCode = []byte("1234")
	err := ValidateInitCode(op, func(f common.Address) (*entrypoint.IStakeManagerDepositInfo, error) {
		return nil, nil
	})

	if err == nil {
		t.Fatalf("got nil, want err")
	}
}

// TestGetStakeFuncReceivesFactory calls checks.ValidateInitCode where initCode exist and calls getStakeFunc
// with the correct factory address.
func TestGetStakeFuncReceivesFactory(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	isCalled := false
	_ = ValidateInitCode(op, func(f common.Address) (*entrypoint.IStakeManagerDepositInfo, error) {
		if f != op.GetFactory() {
			t.Fatalf("got %s, want %s", f.String(), op.GetFactory())
		}

		isCalled = true
		return nil, nil
	})

	if !isCalled {
		t.Fatalf("getStakeFunc was not called")
	}
}

// TestInitCodeExists calls checks.ValidateInitCode where valid initCode does exist. Expect nil.
func TestInitCodeExists(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	err := ValidateInitCode(op, func(f common.Address) (*entrypoint.IStakeManagerDepositInfo, error) {
		return nil, nil
	})

	if err != nil {
		t.Fatalf(`got err %v, want nil`, err)
	}
}
