package state

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestWithZeroAddressOverride(t *testing.T) {
	os, err := ParseOverrideData(map[string]any{})
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if _, ok := os[common.HexToAddress("0x")]; ok {
		t.Fatal("New OverrideSet should not contain zero address override")
	}

	os = WithZeroAddressOverride(os)
	if oa, ok := os[common.HexToAddress("0x")]; !ok {
		t.Fatal("OverrideSet does not contain OverrideAccount")
	} else if oa.Balance.ToInt().String() != maxUint96.String() {
		t.Fatalf("got %s, want %s", oa.Balance.ToInt().String(), maxUint96.String())
	}
}
