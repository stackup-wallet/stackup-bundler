package gas

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/execution"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// CallGasEstimate uses the simulateHandleOp method on the EntryPoint to derive an estimate for callGasLimit.
//
// TODO: This function requires an eth_call and a debug_traceCall. It could probably be optimized further by
// just using a debug_traceCall.
func CallGasEstimate(
	rpc *rpc.Client,
	from common.Address,
	op *userop.UserOperation,
	chainID *big.Int,
	tracer string,
) (uint64, error) {
	data, err := op.ToMap()
	if err != nil {
		return 0, err
	}

	// Set MaxPriorityFeePerGas = MaxFeePerGas to simplify callGasLimit calculation.
	data["maxPriorityFeePerGas"] = hexutil.EncodeBig(op.MaxFeePerGas)
	simOp, err := userop.New(data)
	if err != nil {
		return 0, err
	}

	sim, err := execution.SimulateHandleOp(rpc, from, simOp, common.Address{}, []byte{})
	if err != nil {
		return 0, err
	}

	if err := execution.TraceSimulateHandleOp(rpc, from, op, chainID, tracer, common.Address{}, []byte{}); err != nil {
		return 0, err
	}

	ov := NewDefaultOverhead()
	tg := big.NewInt(0).Div(sim.Paid, op.MaxFeePerGas)
	cgl := big.NewInt(0).Add(big.NewInt(0).Sub(tg, sim.PreOpGas), big.NewInt(int64(ov.fixed)))
	min := ov.NonZeroValueCall()
	if cgl.Cmp(min) >= 1 {
		return cgl.Uint64(), nil
	}
	return min.Uint64(), nil
}
