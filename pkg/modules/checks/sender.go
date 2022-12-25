package checks

import (
	"errors"

	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// ValidateSender accepts a userOp and a generic function that can retrieve the bytecode of the sender.
// Either the sender is deployed (non-zero length bytecode) or the initCode is not empty (but not both).
func ValidateSender(op *userop.UserOperation, gc GetCodeFunc) error {
	bytecode, err := gc(op.Sender)
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
