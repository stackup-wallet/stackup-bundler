package checks

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// TestCheckSender calls checks.CheckSender where sender is existing contract and initCode is empty
func TestCheckSender(t *testing.T) {
	server := testutils.EthMock(testutils.MethodMocks{
		"eth_getCode": "0x1234",
	})
	defer server.Close()

	eth, _ := ethclient.Dial(server.URL)
	op := &userop.UserOperation{
		Sender:   common.HexToAddress("0x"),
		InitCode: common.Hex2Bytes("0x"),
	}
	if err := checkSender(eth, op); err != nil {
		t.Fatalf(`got err %v, want nil`, err)
	}
}
