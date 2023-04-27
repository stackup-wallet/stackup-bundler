package gas

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/execution"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// EstimateGas uses the simulateHandleOp method on the EntryPoint to derive an estimate for
// verificationGasLimit and callGasLimit.
//
// TODO: This function requires an eth_call and a debug_traceCall. It could probably be optimized further by
// just using a debug_traceCall.
func EstimateGas(
	rpc *rpc.Client,
	from common.Address,
	op *userop.UserOperation,
	ov *Overhead,
	chainID *big.Int,
	tracer string,
) (verificationGas uint64, callGas uint64, err error) {
	if op.MaxFeePerGas.Cmp(big.NewInt(0)) != 1 {
		return 0, 0, errors.NewRPCError(
			errors.INVALID_FIELDS,
			"maxFeePerGas must be more than 0",
			nil,
		)
	}
	data, err := op.ToMap()
	if err != nil {
		return 0, 0, err
	}

	// Set MaxPriorityFeePerGas = MaxFeePerGas to simplify callGasLimit calculation.
	data["maxPriorityFeePerGas"] = hexutil.EncodeBig(op.MaxFeePerGas)
	simOp, err := userop.New(data)
	if err != nil {
		return 0, 0, err
	}

	sim, err := execution.SimulateHandleOp(rpc, from, simOp, common.Address{}, []byte{})
	if err != nil {
		return 0, 0, err
	}

	if err := execution.TraceSimulateHandleOp(rpc, from, op, chainID, tracer, common.Address{}, []byte{}); err != nil {
		return 0, 0, err
	}

	tg := big.NewInt(0).Div(sim.Paid, op.MaxFeePerGas)
	cgl := big.NewInt(0).Add(big.NewInt(0).Sub(tg, sim.PreOpGas), big.NewInt(int64(ov.fixed)))
	min := ov.NonZeroValueCall()
	if cgl.Cmp(min) >= 1 {
		return sim.PreOpGas.Uint64(), cgl.Uint64(), nil
	}
	return sim.PreOpGas.Uint64(), min.Uint64(), nil
}
