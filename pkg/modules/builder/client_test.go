package builder

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/metachris/flashbotsrpc"
	"github.com/stackup-wallet/stackup-bundler/internal/testutils"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func TestSendUserOperationWithAllUpstreamErrors(t *testing.T) {
	n := testutils.RpcMock(testutils.MethodMocks{
		"eth_blockNumber":           "0x1",
		"eth_gasPrice":              "0x1",
		"eth_getTransactionCount":   "0x1",
		"eth_estimateGas":           "0x1",
		"eth_getBlockByNumber":      testutils.NewBlockMock(),
		"eth_getTransactionReceipt": testutils.NewTransactionReceiptMock(),
	})
	r, _ := rpc.Dial(n.URL)
	eth := ethclient.NewClient(r)

	bb1 := testutils.BadBuilderRpcMock()
	bb2 := testutils.BadBuilderRpcMock()
	fb := flashbotsrpc.NewBuilderBroadcastRPC([]string{bb1.URL, bb2.URL})
	fn := New(testutils.DummyEOA, eth, fb, testutils.DummyEOA.Address, 1).SendUserOperation()

	if err := fn(
		modules.NewBatchHandlerContext(
			[]*userop.UserOperationV06{testutils.MockValidInitV06UserOp()},
			common.HexToAddress("0x"),
			testutils.ChainID,
			big.NewInt(1),
			big.NewInt(1),
			big.NewInt(1),
		),
	); !errors.Is(err, ErrFlashbotsBroadcastBundle) {
		t.Fatalf("got %v, want ErrFlashbotsBroadcastBundle", err)
	}
}

func TestSendUserOperationWithPartialUpstreamErrors(t *testing.T) {
	n := testutils.RpcMock(testutils.MethodMocks{
		"eth_blockNumber":           "0x1",
		"eth_gasPrice":              "0x1",
		"eth_getTransactionCount":   "0x1",
		"eth_estimateGas":           "0x1",
		"eth_getBlockByNumber":      testutils.NewBlockMock(),
		"eth_getTransactionReceipt": testutils.NewTransactionReceiptMock(),
	})
	r, _ := rpc.Dial(n.URL)
	eth := ethclient.NewClient(r)

	bb1 := testutils.RpcMock(testutils.MethodMocks{
		"eth_sendBundle": map[string]string{
			"bundleHash": testutils.MockHash,
		},
	})
	bb2 := testutils.BadBuilderRpcMock()
	fb := flashbotsrpc.NewBuilderBroadcastRPC([]string{bb1.URL, bb2.URL})
	fn := New(testutils.DummyEOA, eth, fb, testutils.DummyEOA.Address, 1).SendUserOperation()

	if err := fn(
		modules.NewBatchHandlerContext(
			[]*userop.UserOperationV06{testutils.MockValidInitV06UserOp()},
			common.HexToAddress("0x"),
			testutils.ChainID,
			big.NewInt(1),
			big.NewInt(1),
			big.NewInt(1),
		),
	); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

func TestSendUserOperationWithNoUpstreamErrors(t *testing.T) {
	n := testutils.RpcMock(testutils.MethodMocks{
		"eth_blockNumber":           "0x1",
		"eth_gasPrice":              "0x1",
		"eth_getTransactionCount":   "0x1",
		"eth_estimateGas":           "0x1",
		"eth_getBlockByNumber":      testutils.NewBlockMock(),
		"eth_getTransactionReceipt": testutils.NewTransactionReceiptMock(),
	})
	r, _ := rpc.Dial(n.URL)
	eth := ethclient.NewClient(r)

	bb1 := testutils.RpcMock(testutils.MethodMocks{
		"eth_sendBundle": map[string]string{
			"bundleHash": testutils.MockHash,
		},
	})
	bb2 := testutils.RpcMock(testutils.MethodMocks{
		"eth_sendBundle": map[string]string{
			"bundleHash": testutils.MockHash,
		},
	})
	fb := flashbotsrpc.NewBuilderBroadcastRPC([]string{bb1.URL, bb2.URL})
	fn := New(testutils.DummyEOA, eth, fb, testutils.DummyEOA.Address, 1).SendUserOperation()

	if err := fn(
		modules.NewBatchHandlerContext(
			[]*userop.UserOperationV06{testutils.MockValidInitV06UserOp()},
			common.HexToAddress("0x"),
			testutils.ChainID,
			big.NewInt(1),
			big.NewInt(1),
			big.NewInt(1),
		),
	); err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}
