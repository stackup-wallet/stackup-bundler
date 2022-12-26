package checks

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// ValidatePaymasterAndData checks the paymasterAndData is either zero bytes or the first 20 bytes contain an
// address that
//
//  1. currently has nonempty code on chain
//  2. has a sufficient deposit to pay for the UserOperation
func ValidatePaymasterAndData(op *userop.UserOperation, gc GetCodeFunc, gs GetStakeFunc) error {
	if len(op.PaymasterAndData) == 0 {
		return nil
	}

	if len(op.PaymasterAndData) < common.AddressLength {
		return errors.New("PaymasterAndData: invalid length")
	}

	pm := op.GetPaymaster()
	bytecode, err := gc(pm)
	if err != nil {
		return err
	}
	if len(bytecode) == 0 {
		return errors.New("paymaster: code not deployed")
	}

	dep, err := gs(pm)
	if err != nil {
		return err
	}
	if dep.Deposit.Cmp(op.GetMaxPrefund()) < 0 {
		return errors.New("paymaster: not enough deposit to cover max prefund")
	}

	return nil
}
