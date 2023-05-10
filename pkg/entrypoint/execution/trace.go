package execution

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	ethRpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/reverts"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/utils"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func TraceSimulateHandleOp(
	rpc *ethRpc.Client,
	entryPoint common.Address,
	op *userop.UserOperation,
	chainID *big.Int,
	customTracer string,
	target common.Address,
	data []byte,
) (*reverts.ExecutionResultRevert, error) {
	ep, err := entrypoint.NewEntrypoint(entryPoint, ethclient.NewClient(rpc))
	if err != nil {
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(utils.DummyPk, chainID)
	if err != nil {
		return nil, err
	}
	auth.GasLimit = math.MaxUint64
	auth.NoSend = true
	tx, err := ep.SimulateHandleOp(auth, entrypoint.UserOperation(*op), target, data)
	if err != nil {
		return nil, err
	}

	var res tracer.BundlerErrorReturn
	req := utils.TraceCallReq{
		From: common.HexToAddress("0x"),
		To:   entryPoint,
		Data: tx.Data(),
	}
	opts := utils.TraceCallOpts{
		Tracer: customTracer,
	}
	if err := rpc.CallContext(context.Background(), &res, "debug_traceCall", &req, "latest", &opts); err != nil {
		return nil, err
	}
	outErr, err := errors.ParseHexToRpcDataError(res.Output)
	if err != nil {
		return nil, err
	}

	sim, simErr := reverts.NewExecutionResult(outErr)
	if simErr != nil {
		fo, foErr := reverts.NewFailedOp(outErr)
		if foErr != nil {
			return nil, fmt.Errorf("%s, %s", simErr, foErr)
		}
		return nil, errors.NewRPCError(errors.REJECTED_BY_EP_OR_ACCOUNT, fo.Reason, fo)
	}

	if len(res.Reverts) != 0 {
		data, err := hexutil.Decode(res.Reverts[len(res.Reverts)-1])
		if err != nil {
			return sim, err
		}

		if len(data) == 0 {
			return sim, errors.NewRPCError(errors.EXECUTION_REVERTED, "execution reverted", nil)
		}

		reason, err := errors.DecodeRevert(data)
		if err != nil {
			return sim, err
		}
		return sim, errors.NewRPCError(errors.EXECUTION_REVERTED, reason, reason)
	}
	return sim, nil
}
