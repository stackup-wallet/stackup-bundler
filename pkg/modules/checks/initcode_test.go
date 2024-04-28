package checks

import (
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

func TestInitCodeDNE(t *testing.T) {
	op := testutils.MockValidInitV06UserOp()
	op.InitCode = []byte{}
	err := ValidateInitCode(op)

	if err != nil {
		t.Fatalf(`got err %v, want nil`, err)
	}
}

func TestInitCodeContainsBadAddress(t *testing.T) {
	op := testutils.MockValidInitV06UserOp()
	op.InitCode = []byte("1234")
	err := ValidateInitCode(op)

	if err == nil {
		t.Fatalf("got nil, want err")
	}
}

func TestInitCodeExists(t *testing.T) {
	op := testutils.MockValidInitV06UserOp()
	err := ValidateInitCode(op)

	if err != nil {
		t.Fatalf(`got err %v, want nil`, err)
	}
}
