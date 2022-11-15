package base

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Checks that the sender is an existing contract, or the initCode is not empty (but not both)
func checkSender(eth *ethclient.Client, op *userop.UserOperation) error {
	bytecode, err := eth.CodeAt(context.Background(), op.Sender, nil)
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
func checkVerificationGas(maxVerificationGas *big.Int, op *userop.UserOperation) error {
	if op.VerificationGasLimit.Cmp(maxVerificationGas) > 0 {
		return fmt.Errorf("verificationGasLimit: exceeds maxVerificationGas of %s", maxVerificationGas.String())
	}

	ov := gas.NewDefaultOverhead()
	pvg := ov.CalcPreVerificationGas(op)
	if op.PreVerificationGas.Cmp(pvg) < 0 {
		return fmt.Errorf("preVerificationGas: below expected gas of %s", pvg.String())
	}

	return nil
}

// Checks the paymasterAndData is either zero bytes or the first 20 bytes contain an address that
//
//  1. is not the zero address
//  2. currently has nonempty code on chain
//  3. has registered and staked
//  4. has a sufficient deposit to pay for the UserOperation
//  5. is not currently banned
func checkPaymasterAndData(eth *ethclient.Client, op *userop.UserOperation, ep *entrypoint.Entrypoint) error {
	if len(op.PaymasterAndData) == 0 {
		return nil
	}

	address := common.BytesToAddress(op.PaymasterAndData)
	if address == common.HexToAddress("0x") {
		return errors.New("paymaster: cannot be the zero address")
	}

	bytecode, err := eth.CodeAt(context.Background(), address, nil)
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
func checkCallGasLimit(eth *ethclient.Client, op *userop.UserOperation) error {
	// TODO: Add implementation
	return nil
}

// The maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client
// is willing to accept. At the minimum, they are sufficiently high to be included with the current
// block.basefee.
func checkFeePerGas(eth *ethclient.Client, op *userop.UserOperation) error {
	// TODO: Add implementation
	return nil
}

// The sender can only have one UserOperation in the mempool. However it can be replaced if
//
//	(i) the nonce remains the same
//	(ii) the new maxPriorityFeePerGas is higher
//	(iii) the new maxFeePerGas is increased equally
func checkDuplicates(mem *mempool.Interface, op *userop.UserOperation, ep common.Address) error {
	op, err := mem.GetOp(ep, op.Sender)
	if err != nil {
		return err
	}
	if op == nil {
		return nil
	}

	if op.Nonce.Cmp(op.Nonce) != 0 {
		return errors.New("sender: Has userOp in mempool with a different nonce")
	}

	if op.MaxPriorityFeePerGas.Cmp(op.MaxPriorityFeePerGas) <= 0 {
		return errors.New("sender: Has userOp in mempool with same or higher priority fee")
	}

	diff := big.NewInt(0)
	mf := big.NewInt(0)
	diff.Sub(op.MaxPriorityFeePerGas, op.MaxPriorityFeePerGas)
	if op.MaxFeePerGas.Cmp(mf.Add(op.MaxFeePerGas, diff)) != 0 {
		return errors.New("sender: Replaced userOp must have an equally higher max fee")
	}

	return nil
}
