package standalone

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

	pvg := gas.NewDefaultOverhead().CalcPreVerificationGas(op)
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
func checkPaymasterAndData(eth *ethclient.Client, op *userop.UserOperation, ep *entrypoint.Entrypoint) error {
	if len(op.PaymasterAndData) == 0 {
		return nil
	}

	if len(op.PaymasterAndData) < common.AddressLength {
		return errors.New("PaymasterAndData: invalid length")
	}

	address := op.GetPaymaster()
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

	// TODO: Implement condition (iv)
	return nil
}

// Checks the callGasLimit is at least the cost of a CALL with non-zero value.
func checkCallGasLimit(op *userop.UserOperation) error {
	cg := gas.NewDefaultOverhead().CalcCallGasCost()
	if op.CallGasLimit.Cmp(cg) < 0 {
		return fmt.Errorf("callGasLimit: below expected gas of %s", cg.String())
	}

	return nil
}

// The maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client
// is willing to accept. At the minimum, they are sufficiently high to be included with the current
// block.basefee.
func checkFeePerGas(eth *ethclient.Client, op *userop.UserOperation) error {
	tip, err := eth.SuggestGasTipCap(context.Background())
	if err != nil {
		return err
	}

	if op.MaxPriorityFeePerGas.Cmp(tip) < 0 {
		return fmt.Errorf("maxPriorityFeePerGas: below expected wei of %s", tip.String())
	}
	if op.MaxFeePerGas.Cmp(op.MaxPriorityFeePerGas) < 0 {
		return fmt.Errorf("maxFeePerGas: must be equal to or greater than maxPriorityFeePerGas")
	}

	return nil
}
