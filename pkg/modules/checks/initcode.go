package checks

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// ValidateInitCode checks if initCode is not empty and gets the factory address. If factory address is valid
// it calls a generic function that can retrieve the stake from the EntryPoint.
func ValidateInitCode(op *userop.UserOperation, gs GetStakeFunc) error {
	if len(op.InitCode) == 0 {
		return nil
	}

	f := op.GetFactory()
	if f == common.HexToAddress("0x") {
		return errors.New("initCode: does not contain a valid address")
	}

	_, err := gs(f)
	if err != nil {
		return err
	}

	return nil
}
