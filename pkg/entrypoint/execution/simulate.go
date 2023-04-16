package execution

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/reverts"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func SimulateHandleOp(
	eth *ethclient.Client,
	entryPoint common.Address,
	op *userop.UserOperation,
	target common.Address,
	data []byte,
) (*reverts.ExecutionResultRevert, error) {
	ep, err := entrypoint.NewEntrypoint(entryPoint, eth)
	if err != nil {
		return nil, err
	}

	rawCaller := &entrypoint.EntrypointRaw{Contract: ep}
	err = rawCaller.Call(
		nil,
		nil,
		"simulateHandleOp",
		entrypoint.UserOperation(*op),
		target,
		data,
	)

	sim, simErr := reverts.NewExecutionResult(err)
	if simErr != nil {
		fo, foErr := reverts.NewFailedOp(err)
		if foErr != nil {
			return nil, fmt.Errorf("%s, %s", simErr, foErr)
		}
		return nil, errors.NewRPCError(errors.REJECTED_BY_EP_OR_ACCOUNT, fo.Reason, fo)
	}

	return sim, nil
}
