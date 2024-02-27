package state

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// TestParseOverrideDataBadKey validates that ParseOverrideData returns an error if input contains a key that
// cannot be decoded into an address.
func TestParseOverrideDataBadKey(t *testing.T) {
	data := map[string]any{
		"NOT_AN_ADDRESS": nil,
	}

	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadKey) {
		t.Fatal("got nil, want ErrBadKey")
	}
}

// TestParseOverrideDataBadValue validates that parseOverrideData returns an error if input contains a value
// of the wrong type.
func TestParseOverrideDataBadValue(t *testing.T) {
	data := map[string]any{
		"0x0000000000000000000000000000000000000000": 1,
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data["0x0000000000000000000000000000000000000000"] = "string"
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data["0x0000000000000000000000000000000000000000"] = []any{}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}
}

// TestParseOverrideDataEmpty validates that parseOverrideData returns the correct OverrideSet if input
// contains an empty map.
func TestParseOverrideDataEmpty(t *testing.T) {
	data := map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{},
	}
	if os, err := ParseOverrideData(data); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if oa, ok := os[common.HexToAddress("0x")]; !ok {
		t.Fatal("OverrideSet does not contain OverrideAccount")
	} else if oa.Nonce != nil ||
		oa.Code != nil ||
		oa.Balance != nil ||
		oa.State != nil ||
		oa.StateDiff != nil {
		t.Fatal("OverrideAccount contains non nil values")
	}
}

// TestParseOverrideDataNonce validates that parseOverrideData returns the correct response if input
// contains a nonce override.
func TestParseOverrideDataNonce(t *testing.T) {
	nonceOv := "0x1"
	data := map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"nonce": nonceOv,
		},
	}
	if os, err := ParseOverrideData(data); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if oa, ok := os[common.HexToAddress("0x")]; !ok {
		t.Fatal("OverrideSet does not contain OverrideAccount")
	} else if oa.Code != nil ||
		oa.Balance != nil ||
		oa.State != nil ||
		oa.StateDiff != nil {
		t.Fatal("OverrideAccount unset fields contains non nil values")
	} else if oa.Nonce.String() != nonceOv {
		t.Fatalf("got %s, want %s", oa.Nonce.String(), nonceOv)
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"nonce": 1,
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"nonce": "1",
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"nonce": []any{},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}
}

// TestParseOverrideDataCode validates that parseOverrideData returns the correct response if input
// contains a code override.
func TestParseOverrideDataCode(t *testing.T) {
	codeOv := "0x01"
	data := map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"code": codeOv,
		},
	}
	if os, err := ParseOverrideData(data); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if oa, ok := os[common.HexToAddress("0x")]; !ok {
		t.Fatal("OverrideSet does not contain OverrideAccount")
	} else if oa.Nonce != nil ||
		oa.Balance != nil ||
		oa.State != nil ||
		oa.StateDiff != nil {
		t.Fatal("OverrideAccount unset fields contains non nil values")
	} else if oa.Code.String() != codeOv {
		t.Fatalf("got %s, want %s", oa.Code.String(), codeOv)
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"code": 1,
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"code": "1",
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"code": []any{},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"code": "0x1",
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}
}

// TestParseOverrideDataBalance validates that parseOverrideData returns the correct response if input
// contains a balance override.
func TestParseOverrideDataBalance(t *testing.T) {
	balOv := "0x1"
	data := map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"balance": balOv,
		},
	}
	if os, err := ParseOverrideData(data); err != nil {
		t.Fatalf("got %v, want nil", err)
	} else if oa, ok := os[common.HexToAddress("0x")]; !ok {
		t.Fatal("OverrideSet does not contain OverrideAccount")
	} else if oa.Nonce != nil ||
		oa.Code != nil ||
		oa.State != nil ||
		oa.StateDiff != nil {
		t.Fatal("OverrideAccount unset fields contains non nil values")
	} else if oa.Balance.String() != balOv {
		t.Fatalf("got %s, want %s", oa.Balance.String(), balOv)
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"balance": 1,
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"balance": "1",
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"balance": []any{},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}
}

// TestParseOverrideDataState validates that parseOverrideData returns the correct response if input
// contains a state override.
func TestParseOverrideDataState(t *testing.T) {
	stateOvKey := common.HexToHash("0xdead")
	stateOvVal := common.HexToHash("0xbeef")
	data := map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"state": map[string]any{
				stateOvKey.String(): stateOvVal.String(),
			},
		},
	}
	os, err := ParseOverrideData(data)
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	oa, ok := os[common.HexToAddress("0x")]
	if !ok {
		t.Fatal("OverrideSet does not contain OverrideAccount")
	} else if oa.Nonce != nil ||
		oa.Code != nil ||
		oa.Balance != nil ||
		oa.StateDiff != nil {
		t.Fatal("OverrideAccount unset fields contains non nil values")
	}
	s := *oa.State
	if s[stateOvKey].String() != stateOvVal.String() {
		t.Fatalf("got %s, want %s", s[stateOvKey].String(), stateOvVal.String())
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"state": 1,
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"state": "1",
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"state": []any{},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"state": map[string]any{
				"1": stateOvVal.String(),
			},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"state": map[string]any{
				"0x1": stateOvVal.String(),
			},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"state": map[string]any{
				stateOvKey.String(): 1,
			},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"state": map[string]any{
				stateOvKey.String(): "1",
			},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"state": map[string]any{
				stateOvKey.String(): []any{},
			},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"state": map[string]any{
				stateOvKey.String(): "0x1",
			},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}
}

// TestParseOverrideDataStateDiff validates that parseOverrideData returns the correct response if input
// contains a stateDiff override.
func TestParseOverrideDataStateDiff(t *testing.T) {
	stateOvKey := common.HexToHash("0xdead")
	stateOvVal := common.HexToHash("0xbeef")
	data := map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"stateDiff": map[string]any{
				stateOvKey.String(): stateOvVal.String(),
			},
		},
	}
	os, err := ParseOverrideData(data)
	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
	oa, ok := os[common.HexToAddress("0x")]
	if !ok {
		t.Fatal("OverrideSet does not contain OverrideAccount")
	} else if oa.Nonce != nil ||
		oa.Code != nil ||
		oa.Balance != nil ||
		oa.State != nil {
		t.Fatal("OverrideAccount unset fields contains non nil values")
	}
	s := *oa.StateDiff
	if s[stateOvKey].String() != stateOvVal.String() {
		t.Fatalf("got %s, want %s", s[stateOvKey].String(), stateOvVal.String())
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"stateDiff": 1,
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"stateDiff": "1",
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"stateDiff": []any{},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"stateDiff": map[string]any{
				"1": stateOvVal.String(),
			},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"stateDiff": map[string]any{
				"0x1": stateOvVal.String(),
			},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"stateDiff": map[string]any{
				stateOvKey.String(): 1,
			},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"stateDiff": map[string]any{
				stateOvKey.String(): "1",
			},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"stateDiff": map[string]any{
				stateOvKey.String(): []any{},
			},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}

	data = map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
			"stateDiff": map[string]any{
				stateOvKey.String(): "0x1",
			},
		},
	}
	if _, err := ParseOverrideData(data); !errors.Is(err, ErrBadValue) {
		t.Fatal("got nil, want ErrBadValue")
	}
}

// TestParseOverrideDataOk validates that parseOverrideData returns the correct OverrideSet if input contains
// a valid map.
func TestParseOverrideDataOk(t *testing.T) {
	stateOvKey := common.HexToHash("0xdead")
	stateOvVal := common.HexToHash("0xbeef")
	data := map[string]any{
		"0x0000000000000000000000000000000000000000": map[string]any{
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
	} else if oa, ok := os[common.HexToAddress("0x")]; !ok {
		t.Fatal("OverrideSet does not contain OverrideAccount")
	} else if oa.Nonce == nil ||
		oa.Code == nil ||
		oa.Balance == nil ||
		oa.State == nil ||
		oa.StateDiff == nil {
		t.Fatal("OverrideAccount contains nil values")
	}
}
