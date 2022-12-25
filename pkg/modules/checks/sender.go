package checks

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
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
