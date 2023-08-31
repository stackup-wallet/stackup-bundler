package execution

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/reverts"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type SimulateInput struct {
	Rpc        *rpc.Client
	EntryPoint common.Address
	Op         *userop.UserOperation

	// Optional params for simulateHandleOps
	Target common.Address
	Data   []byte
}

func SimulateHandleOp(in *SimulateInput) (*reverts.ExecutionResultRevert, error) {
	ep, err := entrypoint.NewEntrypoint(in.EntryPoint, ethclient.NewClient(in.Rpc))
	if err != nil {
		return nil, err
	}

	rawCaller := &entrypoint.EntrypointRaw{Contract: ep}
	err = rawCaller.Call(
		nil,
		nil,
		"simulateHandleOp",
		entrypoint.UserOperation(*in.Op),
		in.Target,
		in.Data,
	)

	sim, simErr := reverts.NewExecutionResult(err)
	if simErr != nil {
		fo, foErr := reverts.NewFailedOp(err)
		if foErr != nil {
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("%s, %s", simErr, foErr)
		}
		return nil, errors.NewRPCError(errors.REJECTED_BY_EP_OR_ACCOUNT, fo.Reason, fo)
	}

	return sim, nil
}
