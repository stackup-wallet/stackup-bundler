package state

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestCopyNil(t *testing.T) {
	if os, err := Copy(nil); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if os == nil {
		t.Fatal("got nil return value")
	} else if len(os) != 0 {
		t.Fatalf("got length %v, want 0", len(os))
	}
}

func TestCopyNewRef(t *testing.T) {
	acc := common.HexToAddress("0x0000000000000000000000000000000000000000")
	stateOvKey := common.HexToHash("0xdead")
	stateOvVal := common.HexToHash("0xbeef")
	data := map[string]any{
		acc.Hex(): map[string]any{
			"nonce":   "0x1",
			"code":    "0x01",
			"balance": "0x1",
			"state": map[string]any{
				stateOvKey.String(): stateOvVal.String(),
			},
			"stateDiff": map[string]any{
				stateOvKey.String(): stateOvVal.String(),
			},
		},
	}
	if os, err := ParseOverrideData(data); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if osCpy, cpyErr := Copy(os); err != nil {
		t.Fatalf("got %v, want nil", cpyErr)
	} else if os[acc].Nonce == osCpy[acc].Nonce ||
		os[acc].Code == osCpy[acc].Code ||
		os[acc].Balance == osCpy[acc].Balance ||
		os[acc].State == osCpy[acc].State ||
		os[acc].StateDiff == osCpy[acc].StateDiff {
		t.Fatal("OverrideAccount contains identical references")
	}
}

func TestCopySameVal(t *testing.T) {
	acc := common.HexToAddress("0x0000000000000000000000000000000000000000")
	stateOvKey := common.HexToHash("0xdead")
	stateOvVal := common.HexToHash("0xbeef")
	data := map[string]any{
		acc.Hex(): map[string]any{
			"nonce":   "0x1",
			"code":    "0x01",
			"balance": "0x1",
			"state": map[string]any{
				stateOvKey.String(): stateOvVal.String(),
			},
			"stateDiff": map[string]any{
				stateOvKey.String(): stateOvVal.String(),
			},
		},
	}
	os, err := ParseOverrideData(data)
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	osCpy, cpyErr := Copy(os)
	if cpyErr != nil {
		t.Fatalf("got %v, want nil", cpyErr)
	}

	state := *os[acc].State
	stateCpy := *osCpy[acc].State
	stateDiff := *os[acc].StateDiff
	stateDiffCpy := *osCpy[acc].StateDiff
	if os[acc].Nonce.String() != osCpy[acc].Nonce.String() ||
		os[acc].Code.String() != osCpy[acc].Code.String() ||
		os[acc].Balance.String() != osCpy[acc].Balance.String() ||
		state[stateOvKey] != stateCpy[stateOvKey] ||
		stateDiff[stateOvKey] != stateDiffCpy[stateOvKey] {
		t.Fatal("OverrideAccount contains different values")
	}
}
