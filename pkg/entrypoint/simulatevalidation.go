package entrypoint

import (
	stdError "errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// SimulateValidation makes a static call to Entrypoint.simulateValidation(userop) and returns the
// results without any state changes.
func SimulateValidation(
	rpc *rpc.Client,
	entryPoint common.Address,
	op *userop.UserOperation,
) (*ValidationResultRevert, error) {
	ep, err := NewEntrypoint(entryPoint, ethclient.NewClient(rpc))
	if err != nil {
		return nil, err
	}

	var res []interface{}
	rawCaller := &EntrypointRaw{Contract: ep}
	err = rawCaller.Call(nil, &res, "simulateValidation", UserOperation(*op))
	if err == nil {
		return nil, stdError.New("unexpected result from simulateValidation")
	}

	sim, simErr := newValidationResultRevert(err)
	if simErr != nil {
		fo, foErr := newFailedOpRevert(err)
		if foErr != nil {
			return nil, fmt.Errorf("%s, %s", simErr, foErr)
		}
		return nil, errors.NewRPCError(errors.REJECTED_BY_EP_OR_ACCOUNT, fo.Reason, fo)
	}

	return sim, nil
}
