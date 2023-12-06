package testutils

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func NewBlockMock() map[string]any {
	return map[string]any{
		"parentHash":       MockHash,
		"sha3Uncles":       MockHash,
		"stateRoot":        MockHash,
		"transactionsRoot": MockHash,
		"receiptsRoot":     MockHash,
		"logsBloom":        "0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddead",
		"difficulty":       "0x0",
		"number":           "0x1",
		"gasLimit":         hexutil.EncodeBig(big.NewInt(30000000)),
		"gasUsed":          hexutil.EncodeBig(big.NewInt(5000000)),
		"timestamp":        hexutil.EncodeUint64(uint64(time.Now().Unix())),
		"extraData":        "0x",
	}
}

func NewTransactionReceiptMock() map[string]any {
	return map[string]any{
		"blockHash":         MockHash,
		"blockNumber":       "0x1",
		"cumulativeGasUsed": "0x1",
		"effectiveGasPrice": "0x1",
		"from":              common.HexToAddress("0x").Hex(),
		"gasUsed":           "0x1",
		"logs":              []any{},
		"logsBloom":         "0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddead",
		"status":            "0x1",
		"to":                common.HexToAddress("0x").Hex(),
		"transactionHash":   MockHash,
		"transactionIndex":  "0x1",
		"type":              "0x2",
	}
}
