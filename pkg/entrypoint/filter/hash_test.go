package filter

import (
	"strings"
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

func TestIsValidUserOpHash(t *testing.T) {
	if ok := IsValidUserOpHash(testutils.MockHash); !ok {
		t.Fatalf("%s: got false, want true", testutils.MockHash)
	}

	allNumHash := strings.ReplaceAll(testutils.MockHash, "dead", "0101")
	if ok := IsValidUserOpHash(allNumHash); !ok {
		t.Fatalf("%s: got false, want true", allNumHash)
	}
}

func TestIsValidUserOpHashAllCaps(t *testing.T) {
	hash := strings.ToUpper(testutils.MockHash)
	if ok := IsValidUserOpHash(hash); !ok {
		t.Fatalf("%s: got false, want true", hash)
	}
}

func TestIsValidUserOpHashEmptyString(t *testing.T) {
	hash := ""
	if ok := IsValidUserOpHash(hash); ok {
		t.Fatalf("%s: got true, want false", hash)
	}
}

func TestIsValidUserOpHashEmptyHexString(t *testing.T) {
	hash := "0x"
	if ok := IsValidUserOpHash(hash); ok {
		t.Fatalf("%s: got true, want false", hash)
	}
}

func TestIsValidUserOpHashNoPrefix(t *testing.T) {
	hash := strings.TrimPrefix(testutils.MockHash, "0x")
	if ok := IsValidUserOpHash(hash); ok {
		t.Fatalf("%s: got true, want false", hash)
	}
}

func TestIsValidUserOpHashTooShort(t *testing.T) {
	hash := "0xdead"
	if ok := IsValidUserOpHash(hash); ok {
		t.Fatalf("%s: got true, want false", hash)
	}
}

func TestIsValidUserOpHashTooLong(t *testing.T) {
	hash := testutils.MockHash + "dead"
	if ok := IsValidUserOpHash(hash); ok {
		t.Fatalf("%s: got true, want false", hash)
	}
}

func TestIsValidUserOpHashInvalidChar(t *testing.T) {
	hash := strings.ReplaceAll(testutils.MockHash, "dead", "zzzz")
	if ok := IsValidUserOpHash(hash); ok {
		t.Fatalf("%s: got true, want false", hash)
	}
}
