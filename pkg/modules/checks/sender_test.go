package checks

import (
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
)

// TestSenderExistAndInitCodeDNE calls checks.CheckSender where sender contract exist and initCode does not.
// Expect nil.
func TestSenderExistAndInitCodeDNE(t *testing.T) {
	server := testutils.EthMock(testutils.MethodMocks{
		"eth_getCode": testutils.MockByteCode,
	})
	defer server.Close()

	eth, _ := ethclient.Dial(server.URL)
	op := testutils.MockValidInitUserOp()
	op.InitCode = []byte{}
	if err := checkSender(eth, op); err != nil {
		t.Fatalf(`got err %v, want nil`, err)
	}
}

// TestSenderAndInitCodeExist calls checks.CheckSender where sender contract and initCode exist. Expect
// error.
func TestSenderAndInitCodeExist(t *testing.T) {
	server := testutils.EthMock(testutils.MethodMocks{
		"eth_getCode": testutils.MockByteCode,
	})
	defer server.Close()

	eth, _ := ethclient.Dial(server.URL)
	op := testutils.MockValidInitUserOp()
	if err := checkSender(eth, op); err == nil {
		t.Fatalf(`got nil, want err`)
	}
}

// TestSenderDNEAndInitCodeExist calls checks.CheckSender where sender contract does not exist and initCode
// does. Expect nil.
func TestSenderDNEAndInitCodeExist(t *testing.T) {
	server := testutils.EthMock(testutils.MethodMocks{
		"eth_getCode": "0x",
	})
	defer server.Close()

	eth, _ := ethclient.Dial(server.URL)
	op := testutils.MockValidInitUserOp()
	if err := checkSender(eth, op); err != nil {
		t.Fatalf(`got err %v, want nil`, err)
	}
}

// TestSenderAndInitCodeDNE calls checks.CheckSender where sender contract and initCode does not exist.
// Expect error.
func TestSenderAndInitCodeDNE(t *testing.T) {
	server := testutils.EthMock(testutils.MethodMocks{
		"eth_getCode": "0x",
	})
	defer server.Close()

	eth, _ := ethclient.Dial(server.URL)
	op := testutils.MockValidInitUserOp()
	op.InitCode = []byte{}
	if err := checkSender(eth, op); err == nil {
		t.Fatalf(`got nil, want err`)
	}
}
