package filter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/methods"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type HashLookupResult struct {
	UserOperation   *userop.UserOperation `json:"userOperation"`
	EntryPoint      string                `json:"entryPoint"`
	BlockNumber     *big.Int              `json:"blockNumber"`
	BlockHash       common.Hash           `json:"blockHash"`
	TransactionHash common.Hash           `json:"transactionHash"`
}

// GetUserOperationByHash filters the EntryPoint contract for UserOperationEvents and returns the
// corresponding UserOp from a given userOpHash.
func GetUserOperationByHash(
	eth *ethclient.Client,
	userOpHash string,
	entryPoint common.Address,
	chainID *big.Int,
) (*HashLookupResult, error) {
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

		hex := hexutil.Encode(tx.Data())
		if strings.HasPrefix(hex, methods.HandleOpsSelector) {
			data := common.Hex2Bytes(hex[len(methods.HandleOpsSelector):])
			args, err := methods.HandleOpsMethod.Inputs.Unpack(data)
			if err != nil {
				return nil, err
			}
			if len(args) != 2 {
				return nil, fmt.Errorf(
					"handleOps: invalid input length: expected 2, got %d",
					len(args),
				)
			}

			// TODO: Find better way to convert this
			ops, ok := args[0].([]struct {
				Sender               common.Address `json:"sender"`
				Nonce                *big.Int       `json:"nonce"`
				InitCode             []uint8        `json:"initCode"`
				CallData             []uint8        `json:"callData"`
				CallGasLimit         *big.Int       `json:"callGasLimit"`
				VerificationGasLimit *big.Int       `json:"verificationGasLimit"`
				PreVerificationGas   *big.Int       `json:"preVerificationGas"`
				MaxFeePerGas         *big.Int       `json:"maxFeePerGas"`
				MaxPriorityFeePerGas *big.Int       `json:"maxPriorityFeePerGas"`
				PaymasterAndData     []uint8        `json:"paymasterAndData"`
				Signature            []uint8        `json:"signature"`
			})
			if !ok {
				return nil, errors.New("handleOps: cannot assert type: ops is not of type []struct{...}")
			}

			for _, abiOp := range ops {
				data, err := json.Marshal(abiOp)
				if err != nil {
					return nil, err
				}

				var op userop.UserOperation
				if err = json.Unmarshal(data, &op); err != nil {
					return nil, err
				}

				if op.GetUserOpHash(entryPoint, chainID).String() == userOpHash {
					return &HashLookupResult{
						UserOperation:   &op,
						EntryPoint:      entryPoint.String(),
						BlockNumber:     receipt.BlockNumber,
						BlockHash:       receipt.BlockHash,
						TransactionHash: it.Event.Raw.TxHash,
					}, nil
				}
			}
		}

	}

	//lint:ignore ST1005 This needs to match the bundler test spec.
	return nil, errors.New("Missing/invalid userOpHash")
}
