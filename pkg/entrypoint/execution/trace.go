package execution

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethRpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/reverts"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/utils"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type TraceInput struct {
	Rpc        *ethRpc.Client
	EntryPoint common.Address
	Op         *userop.UserOperation
	ChainID    *big.Int

	// Optional params for simulateHandleOps
	Target      common.Address
	Data        []byte
	TraceFeeCap *big.Int
}

type TraceOutput struct {
	Trace  *tracer.BundlerExecutionReturn
	Result *reverts.ExecutionResultRevert
	Event  *entrypoint.EntrypointUserOperationEvent
}

func parseUserOperationEvent(
	entryPoint common.Address,
	ep *entrypoint.Entrypoint,
	log *tracer.LogInfo,
) (*entrypoint.EntrypointUserOperationEvent, error) {
	if log == nil {
		return nil, nil
	}

	topics := []common.Hash{}
	for _, topic := range log.Topics {
		topics = append(topics, common.HexToHash(topic))
	}
	data, err := hexutil.Decode(log.Data)
	if err != nil {
		return nil, err
	}

	ev, err := ep.ParseUserOperationEvent(types.Log{
		Address: entryPoint,
		Topics:  topics,
		Data:    data,
	})
	if err != nil {
		return nil, err
	}

	return ev, nil
}

func TraceSimulateHandleOp(in *TraceInput) (*TraceOutput, error) {
	ep, err := entrypoint.NewEntrypoint(in.EntryPoint, ethclient.NewClient(in.Rpc))
	if err != nil {
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(utils.DummyPk, in.ChainID)
	if err != nil {
		return nil, err
	}
	auth.GasLimit = math.MaxUint64
	auth.NoSend = true
	mf := in.Op.MaxFeePerGas
	if in.TraceFeeCap != nil {
		mf = in.TraceFeeCap
	}
	tx, err := ep.SimulateHandleOp(auth, entrypoint.UserOperation(*in.Op), in.Target, in.Data)
	if err != nil {
		return nil, err
	}
	out := &TraceOutput{}

	var res tracer.BundlerExecutionReturn
	req := utils.TraceCallReq{
		From:         common.HexToAddress("0x"),
		To:           in.EntryPoint,
		Data:         tx.Data(),
		MaxFeePerGas: hexutil.Big(*mf),
	}
	opts := utils.TraceCallOpts{
		Tracer:         tracer.Loaded.BundlerExecutionTracer,
		StateOverrides: utils.DefaultStateOverrides,
	}
	if err := in.Rpc.CallContext(context.Background(), &res, "debug_traceCall", &req, "latest", &opts); err != nil {
		return nil, err
	}
	outErr, err := errors.ParseHexToRpcDataError(res.Output)
	if err != nil {
		return nil, err
	}
	if res.ValidationOOG {
		return nil, errors.NewRPCError(errors.EXECUTION_REVERTED, "validation OOG", nil)
	}
	out.Trace = &res

	sim, simErr := reverts.NewExecutionResult(outErr)
	if simErr != nil {
		fo, foErr := reverts.NewFailedOp(outErr)
		if foErr != nil && res.Error != "" {
			return nil, errors.NewRPCError(errors.EXECUTION_REVERTED, res.Error, nil)
		} else if foErr != nil {
			return nil, fmt.Errorf("%s, %s", simErr, foErr)
		}
		return nil, errors.NewRPCError(errors.REJECTED_BY_EP_OR_ACCOUNT, fo.Reason, fo)
	}
	out.Result = sim

	if len(res.Reverts) != 0 {
		data, err := hexutil.Decode(res.Reverts[len(res.Reverts)-1])
		if err != nil {
			return out, err
		}

		if len(data) == 0 {
			if res.ExecutionOOG {
				return out, errors.NewRPCError(errors.EXECUTION_REVERTED, "execution OOG", nil)
			}
			return out, errors.NewRPCError(errors.EXECUTION_REVERTED, "execution reverted", nil)
		}

		reason, revErr := errors.DecodeRevert(data)
		if revErr != nil {
			code, panErr := errors.DecodePanic(data)
			if panErr != nil {
				return nil, fmt.Errorf("%s, %s", revErr, panErr)
			}

			return out, errors.NewRPCError(
				errors.EXECUTION_REVERTED,
				fmt.Sprintf("panic encountered: %s", code),
				code,
			)
		}
		return out, errors.NewRPCError(errors.EXECUTION_REVERTED, reason, reason)
	}

	ev, err := parseUserOperationEvent(in.EntryPoint, ep, res.UserOperationEvent)
	if err != nil {
		return out, err
	}
	out.Event = ev

	return out, nil
}
