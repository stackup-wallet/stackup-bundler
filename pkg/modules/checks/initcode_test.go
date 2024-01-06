package checks

import (
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

// TestInitCodeDNE calls checks.ValidateInitCode where initCode does not exist. Expect nil.
func TestInitCodeDNE(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	op.InitCode = []byte{}
	err := ValidateInitCode(op)

	if err != nil {
		t.Fatalf(`got err %v, want nil`, err)
	}
}

// TestInitCodeContainsAddress calls checks.ValidateInitCode where initCode exist without a valid address.
// Expect error.
func TestInitCodeContainsAddress(t *testing.T) {
	op := testutils.MockValidInitUserOp()
	op.InitCode = []byte("1234")
	err := ValidateInitCode(op)

	if err == nil {
		t.Fatalf("got nil, want err")
	}
}
