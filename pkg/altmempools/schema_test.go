package altmempools_test

import (
	"testing"

	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/altmempools"
)

func TestValidatesCompliantAltMempool(t *testing.T) {
	if err := altmempools.Schema.Validate(testutils.AltMempoolMock()); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

func TestValidatesBadForbiddenOpcode(t *testing.T) {
	alt := testutils.AltMempoolMock()
	alt["allowlist"] = []any{
		map[string]any{
			"description": "Mock forbiddenOpcode rule",
			"rule":        "forbiddenOpcode",
		},
	}

	if err := altmempools.Schema.Validate(alt); err == nil {
		t.Fatalf("got nil, want err")
	}
}

func TestValidatesBadForbiddenPrecompile(t *testing.T) {
	alt := testutils.AltMempoolMock()
	alt["allowlist"] = []any{
		map[string]any{
			"description": "Mock forbiddenPrecompile rule",
			"rule":        "forbiddenPrecompile",
		},
	}

	if err := altmempools.Schema.Validate(alt); err == nil {
		t.Fatalf("got nil, want err")
	}
}

func TestValidatesBadInvalidStorageAccess(t *testing.T) {
	alt := testutils.AltMempoolMock()
	alt["allowlist"] = []any{
		map[string]any{
			"description": "Mock invalidStorageAccess rule",
			"rule":        "invalidStorageAccess",
		},
	}

	if err := altmempools.Schema.Validate(alt); err == nil {
		t.Fatalf("got nil, want err")
	}
}

func TestValidatesBadNotStaked(t *testing.T) {
	alt := testutils.AltMempoolMock()
	alt["allowlist"] = []any{
		map[string]any{
			"description": "Mock notStaked rule",
			"rule":        "notStaked",
		},
	}

	if err := altmempools.Schema.Validate(alt); err == nil {
		t.Fatalf("got nil, want err")
	}
}
