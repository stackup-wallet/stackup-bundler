package utils_test

import (
	"math/big"
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/utils"
)

// TestAddBuffer verifies that AddBuffer returns a new big.Int that is increased by a given percent value.
func TestAddBuffer(t *testing.T) {
	factor := int64(10)
	amt := big.NewInt(10)
	want := big.NewInt(11)
	if utils.AddBuffer(amt, factor).Cmp(want) != 0 {
		t.Fatalf("got %d, want %d", utils.AddBuffer(amt, factor).Int64(), want.Int64())
	}
}

// TestAddBufferZeroFactor verifies that AddBuffer returns a new big.Int that is unchanged when the factor is
// 0.
func TestAddBufferZeroFactor(t *testing.T) {
	factor := int64(0)
	amt := big.NewInt(100)
	if utils.AddBuffer(amt, factor).Cmp(amt) != 0 {
		t.Fatalf("got %d, want %d", utils.AddBuffer(amt, factor).Int64(), amt.Int64())
	}
}

// TestAddBufferNilAmt verifies that AddBuffer returns nil if amt is nil.
func TestAddBufferNilAmt(t *testing.T) {
	factor := int64(10)
	if utils.AddBuffer(nil, factor) != nil {
		t.Fatalf("got %d, want nil", utils.AddBuffer(nil, factor).Int64())
	}
}
