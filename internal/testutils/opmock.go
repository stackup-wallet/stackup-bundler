package testutils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

var (
	MockUserOpData = map[string]any{
		"sender":               "0xa13D69573f994bf662C2714560c44dd7266FC547",
		"nonce":                "0x0",
		"initCode":             "0xe19e9755942bb0bd0cccce25b1742596b8a8250b3bf2c3e700000000000000000000000078d4f01f56b982a3b03c4e127a5d3afa8ebee6860000000000000000000000008b388a082f370d8ac2e2b3997e9151168bd09ff50000000000000000000000000000000000000000000000000000000000000000",
		"callData":             "0x80c5c7d0000000000000000000000000a13d69573f994bf662c2714560c44dd7266fc547000000000000000000000000000000000000000000000000016345785d8a000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000000",
		"callGasLimit":         "0x558c",
		"verificationGasLimit": "0x129727",
		"maxFeePerGas":         "0xa862145e",
		"maxPriorityFeePerGas": "0xa8621440",
		"paymasterAndData":     "0x",
		"preVerificationGas":   "0xc650",
		"signature":            "0xa925dcc5e5131636e244d4405334c25f034ebdd85c0cb12e8cdb13c15249c2d466d0bade18e2cafd3513497f7f968dcbb63e519acd9b76dcae7acd61f11aa8421b",
	}
	MockByteCode = common.Hex2Bytes("6080604052")
)

// Returns a valid initial userOperation for an EIP-4337 account.
func MockValidInitUserOp() *userop.UserOperation {
	op, _ := userop.New(MockUserOpData)
	return op
}

func IsOpsEqual(op1 *userop.UserOperation, op2 *userop.UserOperation) bool {
	return cmp.Equal(
		op1,
		op2,
		cmp.Comparer(func(a *big.Int, b *big.Int) bool { return a.Cmp(b) == 0 }),
	)
}

func GetOpsDiff(op1 *userop.UserOperation, op2 *userop.UserOperation) string {
	return cmp.Diff(
		op1,
		op2,
		cmp.Comparer(func(a *big.Int, b *big.Int) bool { return a.Cmp(b) == 0 }),
	)
}
