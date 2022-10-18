package userop

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/internal/web3utils"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
)

type UserOperation struct {
	Sender               string   `json:"sender" mapstructure:"sender" validate:"required,eth_addr"`
	Nonce                *big.Int `json:"nonce" mapstructure:"nonce" validate:"required"`
	InitCode             []byte   `json:"initCode"  mapstructure:"initCode" validate:"required"`
	CallData             []byte   `json:"callData" mapstructure:"callData" validate:"required"`
	CallGasLimit         *big.Int `json:"callGasLimit" mapstructure:"callGasLimit" validate:"required"`
	VerificationGasLimit *big.Int `json:"verificationGasLimit" mapstructure:"verificationGasLimit" validate:"required"`
	PreVerificationGas   *big.Int `json:"preVerificationGas" mapstructure:"preVerificationGas" validate:"required"`
	MaxFeePerGas         *big.Int `json:"maxFeePerGas" mapstructure:"maxFeePerGas" validate:"required"`
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas" mapstructure:"maxPriorityFeePerGas" validate:"required"`
	PaymasterAndData     []byte   `json:"paymasterAndData" mapstructure:"paymasterAndData" validate:"required"`
	Signature            []byte   `json:"signature" mapstructure:"signature" validate:"required"`
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

// Checks that the verificationGasLimit is sufficiently low (<= MAX_VERIFICATION_GAS) and the
// preVerificationGas is sufficiently high (enough to pay for the calldata gas cost of serializing
// the UserOperation plus PRE_VERIFICATION_OVERHEAD_GAS)
func (op *UserOperation) CheckVerificationGasLimits(client *ethclient.Client) error {
	// TODO: Add implementation
	return nil
}

// Checks the paymasterAndData is either zero bytes or the first 20 bytes contain an address that
// (i) is not the zero address,
// (ii) currently has nonempty code on chain,
// (iii) has registered and staked,
// (iv) has a sufficient deposit to pay for the UserOperation,
// and (v) is not currently banned.
func (op *UserOperation) CheckPaymasterAndData(client *ethclient.Client, ep *entrypoint.Entrypoint) error {
	if len(op.PaymasterAndData) == 0 {
		return nil
	}

	address := common.BytesToAddress(op.PaymasterAndData)
	if web3utils.IsZeroAddress(address) {
		return errors.New("paymaster: cannot be the zero address")
	}

	bytecode, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		return err
	}
	if len(bytecode) == 0 {
		return errors.New("paymaster: code not deployed")
	}

	dep, err := ep.GetDepositInfo(&bind.CallOpts{}, address)
	if err != nil {
		return errors.New("paymaster: failed to get deposit info")
	}
	if !dep.Staked {
		return errors.New("paymaster: not staked on the entrypoint")
	}

	// TODO: Implement condition (iv) and (v)

	return nil
}

// Checks the callGasLimit is at least the cost of a CALL with non-zero value.
// See https://github.com/wolflo/evm-opcodes/blob/main/gas.md#aa-1-call
func (op *UserOperation) CheckCallGasLimit(client *ethclient.Client) error {
	// TODO: Add implementation
	return nil
}

// The maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client
// is willing to accept. At the minimum, they are sufficiently high to be included with the current
// block.basefee.
func (op *UserOperation) CheckFeePerGas(client *ethclient.Client) error {
	// TODO: Add implementation
	return nil
}

// The sender doesn’t have another UserOperation already present in the pool (or it replaces an existing
// entry with the same sender and nonce, with a higher maxPriorityFeePerGas and an equally increased
// maxFeePerGas). Only one UserOperation per sender may be included in a single batch.
func (op *UserOperation) CheckDuplicate(client *ethclient.Client) error {
	// TODO: Add implementation
	return nil
}
