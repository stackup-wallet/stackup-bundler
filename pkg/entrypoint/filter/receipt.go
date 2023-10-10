package filter

import (
	"context"
	"errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type parsedTransaction struct {
	BlockHash         common.Hash    `json:"blockHash"`
	BlockNumber       string         `json:"blockNumber"`
	From              common.Address `json:"from"`
	CumulativeGasUsed string         `json:"cumulativeGasUsed"`
	GasUsed           string         `json:"gasUsed"`
	Logs              []*types.Log   `json:"logs"`
	LogsBloom         types.Bloom    `json:"logsBloom"`
	TransactionHash   common.Hash    `json:"transactionHash"`
	TransactionIndex  string         `json:"transactionIndex"`
	EffectiveGasPrice string         `json:"effectiveGasPrice"`
}

type UserOperationReceipt struct {
	UserOpHash    common.Hash        `json:"userOpHash"`
	Sender        common.Address     `json:"sender"`
	Paymaster     common.Address     `json:"paymaster"`
	Nonce         string             `json:"nonce"`
	Success       bool               `json:"success"`
	ActualGasCost string             `json:"actualGasCost"`
	ActualGasUsed string             `json:"actualGasUsed"`
	From          common.Address     `json:"from"`
	Receipt       *parsedTransaction `json:"receipt"`
	Logs          []*types.Log       `json:"logs"`
}

// GetUserOperationReceipt filters the EntryPoint contract for UserOperationEvents and returns a receipt for
// both the UserOperation and accompanying transaction.
func GetUserOperationReceipt(
	eth *ethclient.Client,
	userOpHash string,
	entryPoint common.Address,
) (*UserOperationReceipt, error) {
	it, err := filterUserOperationEvent(eth, userOpHash, entryPoint)
	if err != nil {
		return nil, err
	}

	if it.Next() {
		receipt, err := eth.TransactionReceipt(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return nil, err
		}
		tx, isPending, err := eth.TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return nil, err
		} else if isPending {
			//lint:ignore ST1005 This needs to match the bundler test spec.
			return nil, errors.New("Missing/invalid userOpHash")
		}
		from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			return nil, err
		}

		txnReceipt := &parsedTransaction{
			BlockHash:         receipt.BlockHash,
			BlockNumber:       hexutil.EncodeBig(receipt.BlockNumber),
			From:              from,
			CumulativeGasUsed: hexutil.EncodeBig(big.NewInt(0).SetUint64(receipt.CumulativeGasUsed)),
			GasUsed:           hexutil.EncodeBig(big.NewInt(0).SetUint64(receipt.GasUsed)),
			Logs:              receipt.Logs,
			LogsBloom:         receipt.Bloom,
			TransactionHash:   receipt.TxHash,
			TransactionIndex:  hexutil.EncodeBig(big.NewInt(0).SetUint64(uint64(receipt.TransactionIndex))),
			EffectiveGasPrice: hexutil.EncodeBig(tx.GasPrice()),
		}
		logs, err := _filterLogs(it, receipt.Logs)
		if err != nil {
			return nil, err
		}
		return &UserOperationReceipt{
			UserOpHash:    it.Event.UserOpHash,
			Sender:        it.Event.Sender,
			Paymaster:     it.Event.Paymaster,
			Nonce:         hexutil.EncodeBig(it.Event.Nonce),
			Success:       it.Event.Success,
			ActualGasCost: hexutil.EncodeBig(it.Event.ActualGasCost),
			ActualGasUsed: hexutil.EncodeBig(it.Event.ActualGasUsed),
			From:          from,
			Receipt:       txnReceipt,
			Logs:          logs,
		}, nil
	}

	//lint:ignore ST1005 This needs to match the bundler test spec.
	return nil, errors.New("Missing/invalid userOpHash")
}

func _filterLogs(userOperationEvent *entrypoint.EntrypointUserOperationEventIterator, logs []*types.Log) ([]*types.Log, error) {
	startIndex := -1
	endIndex := -1
	for index := range logs {
		log := logs[index]
		//event BeforeExecution()
		if log.Topics[0] == common.HexToHash("0xbb47ee3e183a558b1a2ff0874b079f3fc5478b7454eacf2bfc5af2ff5878f972") {
			startIndex = index
			endIndex = index
		} else if log.Topics[0] == userOperationEvent.Event.Raw.Topics[0] {
			if log.Topics[1] == userOperationEvent.Event.Raw.Topics[1] {
				endIndex = index
			} else {
				if endIndex == -1 {
					startIndex = index
				}
			}
		}
	}
	if endIndex == -1 {
		return nil, errors.New("fatal: no UserOperationEvent in logs")
	}
	return logs[startIndex+1 : endIndex], nil
}
