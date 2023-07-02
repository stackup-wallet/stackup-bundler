package gas

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/execution"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/reverts"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func isPrefundNotPaid(err error) bool {
	return strings.HasPrefix(err.Error(), "AA21") || strings.HasPrefix(err.Error(), "AA31")
}

func isValidationOOG(err error) bool {
	return strings.HasPrefix(err.Error(), "AA13") || strings.Contains(err.Error(), "validation OOG")
}

func isExecutionOOG(err error) bool {
	return strings.Contains(err.Error(), "execution OOG")
}

func runSimulations(rpc *rpc.Client,
	from common.Address,
	op *userop.UserOperation,
	chainID *big.Int,
	tracer string) (*reverts.ExecutionResultRevert, error) {
	data, err := op.ToMap()
	if err != nil {
		return nil, err
	}

	// Set MaxPriorityFeePerGas = MaxFeePerGas to simplify downstream calculations.
	data["maxPriorityFeePerGas"] = hexutil.EncodeBig(op.MaxFeePerGas)

	// Setting default values for gas limits.
	data["verificationGasLimit"] = hexutil.EncodeBig(big.NewInt(0))
	data["callGasLimit"] = hexutil.EncodeBig(big.NewInt(0))

	// Maintain an outer reference to the latest simulation error.
	var simErr error

	// Find the maximal verificationGasLimit to simulate from.
	l := 0
	r := MaxGasLimit
	for l <= r {
		m := (l + r) / 2

		data["verificationGasLimit"] = hexutil.EncodeBig(big.NewInt(int64(m)))
		simOp, err := userop.New(data)
		if err != nil {
			return nil, err
		}
		sim, err := execution.TraceSimulateHandleOp(
			rpc,
			from,
			simOp,
			chainID,
			tracer,
			common.Address{},
			[]byte{},
		)
		simErr = err
		if err != nil {
			if isPrefundNotPaid(err) {
				// VGL too high, go lower.
				r = m - 1
				continue
			}
			if isValidationOOG(err) {
				// VGL too low, go higher.
				l = m + 1
				continue
			}
			// CGL is set to 0 and execution will always be OOG. Ignore it.
			if !isExecutionOOG(err) {
				return nil, err
			}
		}

		data["verificationGasLimit"] = hexutil.EncodeBig(sim.PreOpGas)
		break
	}
	if simErr != nil && !isExecutionOOG(simErr) {
		return nil, simErr
	}

	// Find the maximal callGasLimit to simulate from.
	l = 0
	r = MaxGasLimit
	var res *reverts.ExecutionResultRevert
	for l <= r {
		m := (l + r) / 2

		data["callGasLimit"] = hexutil.EncodeBig(big.NewInt(int64(m)))
		simOp, err := userop.New(data)
		if err != nil {
			return nil, err
		}
		sim, err := execution.TraceSimulateHandleOp(
			rpc,
			from,
			simOp,
			chainID,
			tracer,
			common.Address{},
			[]byte{},
		)
		simErr = err
		if err != nil {
			if isPrefundNotPaid(err) {
				// CGL too high, go lower.
				r = m - 1
				continue
			}
			if isExecutionOOG(err) {
				// CGL too low, go higher.
				l = m + 1
				continue
			}
			return nil, err
		}

		res = sim
		break
	}
	if simErr != nil {
		return nil, simErr
	}

	return res, nil
}

// EstimateGas uses the simulateHandleOp method on the EntryPoint to derive an estimate for
// verificationGasLimit and callGasLimit.
func EstimateGas(
	rpc *rpc.Client,
	from common.Address,
	op *userop.UserOperation,
	ov *Overhead,
	chainID *big.Int,
	tracer string,
) (verificationGas uint64, callGas uint64, err error) {
	// Skip if maxFeePerGas is zero.
	if op.MaxFeePerGas.Cmp(big.NewInt(0)) != 1 {
		return 0, 0, errors.NewRPCError(
			errors.INVALID_FIELDS,
			"maxFeePerGas must be more than 0",
			nil,
		)
	}

	// Estimate gas limits using a binary search approach.
	sim, err := runSimulations(rpc, from, op, chainID, tracer)
	if err != nil {
		return 0, 0, err
	}

	// Return verificationGasLimit and callGasLimit.
	tg := big.NewInt(0).Div(sim.Paid, op.MaxFeePerGas)
	cgl := big.NewInt(0).Add(big.NewInt(0).Sub(tg, sim.PreOpGas), big.NewInt(int64(ov.intrinsicFixed)))
	min := ov.NonZeroValueCall()
	if cgl.Cmp(min) >= 1 {
		return sim.PreOpGas.Uint64(), cgl.Uint64(), nil
	}
	return sim.PreOpGas.Uint64(), min.Uint64(), nil
}
