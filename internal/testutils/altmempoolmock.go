package testutils

import "github.com/ethereum/go-ethereum/common/hexutil"

func AltMempoolMock() map[string]any {
	return map[string]any{
		"description": "Mock Alt Mempool",
		"chainIds":    []any{hexutil.EncodeBig(ChainID)},
		"allowlist": []any{
			map[string]any{
				"description": "Mock forbiddenOpcode rule",
				"rule":        "forbiddenOpcode",
				"entity":      "account",
				"contract":    "0x0000000000000000000000000000000000000000",
				"opcode":      "GAS",
			},
			map[string]any{
				"description": "Mock forbiddenPrecompile rule",
				"rule":        "forbiddenPrecompile",
				"entity":      "account",
				"contract":    "0x0000000000000000000000000000000000000000",
				"precompile":  "0x0000000000000000000000000000000000000000",
			},
			map[string]any{
				"description": "Mock invalidStorageAccess rule",
				"rule":        "invalidStorageAccess",
				"entity":      "account",
				"contract":    "0x0000000000000000000000000000000000000000",
				"slot":        "0x0000000000000000000000000000000000000000",
			},
			map[string]any{
				"description": "Mock notStaked rule",
				"rule":        "notStaked",
				"entity":      "0x0000000000000000000000000000000000000000",
			},
		},
	}
}
