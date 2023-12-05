package state

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestWithMaxBalanceOverride(t *testing.T) {
	os, err := ParseOverrideData(map[string]any{})
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if _, ok := os[common.HexToAddress("0x")]; ok {
		t.Fatal("New OverrideSet should not contain zero address override")
	}

	os = WithMaxBalanceOverride(common.HexToAddress("0x"), os)
	if oa, ok := os[common.HexToAddress("0x")]; !ok {
		t.Fatal("OverrideSet does not contain OverrideAccount")
	} else if oa.Balance.ToInt().String() != maxUint96.String() {
		t.Fatalf("got %s, want %s", oa.Balance.ToInt().String(), maxUint96.String())
	}
}

func TestWithMaxBalanceOverrideNil(t *testing.T) {
	os := WithMaxBalanceOverride(common.HexToAddress("0x"), nil)
	if oa, ok := os[common.HexToAddress("0x")]; !ok {
		t.Fatal("OverrideSet does not contain OverrideAccount")
	} else if oa.Balance.ToInt().String() != maxUint96.String() {
		t.Fatalf("got %s, want %s", oa.Balance.ToInt().String(), maxUint96.String())
	}
}

func TestWithMaxBalanceOverrideNoop(t *testing.T) {
	bal := big.NewInt(1)
	os := WithMaxBalanceOverride(common.HexToAddress("0x"), OverrideSet{
		common.HexToAddress("0x"): OverrideAccount{
			Balance: (*hexutil.Big)(bal),
		},
	})
	if oa, ok := os[common.HexToAddress("0x")]; !ok {
		t.Fatal("OverrideSet does not contain OverrideAccount")
	} else if oa.Balance.ToInt().String() != bal.String() {
		t.Fatalf("got %s, want %s", oa.Balance.ToInt().String(), bal.String())
	}
}
