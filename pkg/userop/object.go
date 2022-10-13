package userop

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type UserOperation struct {
	Sender               string  `json:"sender" mapstructure:"sender" validate:"required,eth_addr"`
	Nonce                big.Int `json:"nonce" mapstructure:"nonce" validate:"required"`
	InitCode             []byte  `json:"initCode"  mapstructure:"initCode" validate:"required"`
	CallData             []byte  `json:"callData" mapstructure:"callData" validate:"required"`
	CallGasLimit         big.Int `json:"callGasLimit" mapstructure:"callGasLimit" validate:"required"`
	VerificationGasLimit big.Int `json:"verificationGasLimit" mapstructure:"verificationGasLimit" validate:"required"`
	PreVerificationGas   big.Int `json:"preVerificationGas" mapstructure:"preVerificationGas" validate:"required"`
	MaxFeePerGas         big.Int `json:"maxFeePerGas" mapstructure:"maxFeePerGas" validate:"required"`
	MaxPriorityFeePerGas big.Int `json:"maxPriorityFeePerGas" mapstructure:"maxPriorityFeePerGas" validate:"required"`
	PaymasterAndData     []byte  `json:"paymasterAndData" mapstructure:"paymasterAndData" validate:"required"`
	Signature            []byte  `json:"signature" mapstructure:"signature" validate:"required"`
}

// Checks that the sender is an existing contract, or the initCode is not empty (but not both)
func (op *UserOperation) CheckSender(client *ethclient.Client) error {
	address := common.HexToAddress(op.Sender)
	bytecode, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		return err
	}

	if len(bytecode) == 0 && len(op.InitCode) == 0 {
		return errors.New("sender: not deployed, initCode must be set")
	}
	if len(bytecode) > 0 && len(op.InitCode) > 0 {
		return errors.New("sender: already deployed, initCode must be empty")
	}

	return nil
}
