package gas

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/execution"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

var (
	fallBackBinarySearchCutoff = int64(30000)
	maxRetries                 = int64(7)
	baseVGLBuffer              = int64(25)
)

func isPrefundNotPaid(err error) bool {
	return strings.Contains(err.Error(), "AA21 didn't pay prefund") ||
		strings.Contains(err.Error(), "AA31 paymaster deposit too low")
}

func isValidationOOG(err error) bool {
	return strings.Contains(err.Error(), "AA40 over verificationGasLimit") ||
		strings.Contains(err.Error(), "AA41 too little verificationGas") ||
		strings.Contains(err.Error(), "AA51 prefund below actualGasCost") ||
		strings.Contains(err.Error(), "AA13 initCode failed or OOG") ||
		strings.Contains(err.Error(), "AA23 reverted (or OOG)") ||
		strings.Contains(err.Error(), "AA33 reverted (or OOG)") ||
		strings.Contains(err.Error(), "return data out of bounds") ||
		strings.Contains(err.Error(), "validation OOG")
}

func isExecutionOOG(err error) bool {
	return strings.Contains(err.Error(), "execution OOG")
}

func isExecutionReverted(err error) bool {
	return strings.Contains(err.Error(), "execution reverted")
}

type EstimateInput struct {
	Rpc         *rpc.Client
	EntryPoint  common.Address
	Op          *userop.UserOperation
	Ov          *Overhead
	ChainID     *big.Int
	MaxGasLimit *big.Int

	attempts int64
	lastVGL  int64
}

// retryEstimateGas will recursively call estimateGas if execution has caused VGL to be under estimated. This
// can occur for edge cases where a paymaster's postOp > gas required during verification or if verification
// has a dependency on CGL. Reset the estimate with a higher buffer on VGL.
func retryEstimateGas(err error, vgl int64, in *EstimateInput) (uint64, uint64, error) {
	if isValidationOOG(err) && in.attempts < maxRetries {
		return EstimateGas(&EstimateInput{
			Rpc:         in.Rpc,
			EntryPoint:  in.EntryPoint,
			Op:          in.Op,
			Ov:          in.Ov,
			ChainID:     in.ChainID,
			MaxGasLimit: in.MaxGasLimit,
			attempts:    in.attempts + 1,
			lastVGL:     vgl,
		})
	}
	return 0, 0, err
}

// EstimateGas uses the simulateHandleOp method on the EntryPoint to derive an estimate for
// verificationGasLimit and callGasLimit.
func EstimateGas(in *EstimateInput) (verificationGas uint64, callGas uint64, err error) {
	// Skip if maxFeePerGas is zero.
	if in.Op.MaxFeePerGas.Cmp(big.NewInt(0)) != 1 {
		return 0, 0, errors.NewRPCError(
			errors.INVALID_FIELDS,
			"maxFeePerGas must be more than 0",
			nil,
		)
	}

	// Set the initial conditions.
	data, err := in.Op.ToMap()
	if err != nil {
		return 0, 0, err
	}
	data["maxPriorityFeePerGas"] = hexutil.EncodeBig(in.Op.MaxFeePerGas)
	data["verificationGasLimit"] = hexutil.EncodeBig(big.NewInt(0))
	data["callGasLimit"] = hexutil.EncodeBig(big.NewInt(0))

	// Find the optimal verificationGasLimit with binary search. Setting gas price to 0 and maxing out the gas
	// limit here would result in certain code paths not being executed which results in an inaccurate gas
	// estimate.
	l := int64(0)
	r := in.MaxGasLimit.Int64()
	f := in.lastVGL
	var simErr error
	for in.lastVGL == 0 && r-l >= fallBackBinarySearchCutoff {
		m := (l + r) / 2

		data["verificationGasLimit"] = hexutil.EncodeBig(big.NewInt(int64(m)))
		simOp, err := userop.New(data)
		if err != nil {
			return 0, 0, err
		}
		_, err = execution.SimulateHandleOp(&execution.SimulateInput{
			Rpc:        in.Rpc,
			EntryPoint: in.EntryPoint,
			Op:         simOp,
		})
		simErr = err
		if err == nil {
			// VGL too high, go lower.
			r = m - 1
			// Set final.
			f = m
			continue
		} else if isPrefundNotPaid(err) {
			// VGL too high, go lower.
			r = m - 1
			continue
		} else if isValidationOOG(err) {
			// VGL too low, go higher.
			l = m + 1
			continue
		} else {
			return 0, 0, err
		}
	}
	if f == 0 {
		return 0, 0, simErr
	}
	f = (f * (100 + baseVGLBuffer)) / 100
	data["verificationGasLimit"] = hexutil.EncodeBig(big.NewInt(int64(f)))

	// Find the optimal callGasLimit by setting the gas price to 0 and maxing out the gas limit. We don't run
	// into the same restrictions here as we do with verificationGasLimit.
	data["maxFeePerGas"] = hexutil.EncodeBig(big.NewInt(0))
	data["maxPriorityFeePerGas"] = hexutil.EncodeBig(big.NewInt(0))
	data["callGasLimit"] = hexutil.EncodeBig(in.MaxGasLimit)
	simOp, err := userop.New(data)
	if err != nil {
		return 0, 0, err
	}
	out, err := execution.TraceSimulateHandleOp(&execution.TraceInput{
		Rpc:         in.Rpc,
		EntryPoint:  in.EntryPoint,
		Op:          simOp,
		ChainID:     in.ChainID,
		TraceFeeCap: in.Op.MaxFeePerGas,
	})
	if err != nil {
		return retryEstimateGas(err, f, in)
	}

	// Calculate final values for verificationGasLimit and callGasLimit.
	vgl := simOp.VerificationGasLimit
	cgl := big.NewInt(int64(out.Trace.ExecutionGasLimit))
	if cgl.Cmp(in.Ov.NonZeroValueCall()) < 0 {
		cgl = in.Ov.NonZeroValueCall()
	}

	// Run a final simulation to check wether or not value transfers are still okay when factoring in the
	// expected gas cost.
	data["maxFeePerGas"] = hexutil.EncodeBig(in.Op.MaxFeePerGas)
	data["maxPriorityFeePerGas"] = hexutil.EncodeBig(in.Op.MaxFeePerGas)
	data["verificationGasLimit"] = hexutil.EncodeBig(vgl)
	data["callGasLimit"] = hexutil.EncodeBig(cgl)
	simOp, err = userop.New(data)
	if err != nil {
		return 0, 0, err
	}
	_, err = execution.TraceSimulateHandleOp(&execution.TraceInput{
		Rpc:        in.Rpc,
		EntryPoint: in.EntryPoint,
		Op:         simOp,
		ChainID:    in.ChainID,
	})
	if err != nil {
		// Execution is successful but one shot tracing has failed. Fallback to binary search with an
		// efficient range. Hitting this point could mean a contract is passing manual gas limits with a
		// static discount, e.g. sub(gas(), STATIC_DISCOUNT). This is not yet accounted for in the tracer.
		if isExecutionOOG(err) || isExecutionReverted(err) {
			l := cgl.Int64()
			r := in.MaxGasLimit.Int64()
			f := int64(0)
			simErr := err
			for r-l >= fallBackBinarySearchCutoff {
				m := (l + r) / 2

				data["callGasLimit"] = hexutil.EncodeBig(big.NewInt(int64(m)))
				simOp, err := userop.New(data)
				if err != nil {
					return 0, 0, err
				}
				_, err = execution.TraceSimulateHandleOp(&execution.TraceInput{
					Rpc:        in.Rpc,
					EntryPoint: in.EntryPoint,
					Op:         simOp,
					ChainID:    in.ChainID,
				})
				simErr = err
				if err == nil {
					// CGL too high, go lower.
					r = m - 1
					// Set final.
					f = m
					continue
				} else if isPrefundNotPaid(err) {
					// CGL too high, go lower.
					r = m - 1
				} else if isExecutionOOG(err) || isExecutionReverted(err) {
					// CGL too low, go higher.
					l = m + 1
					continue
				} else {
					// Unexpected error.
					return 0, 0, err
				}
			}
			if f == 0 {
				return 0, 0, simErr
			}
			return simOp.VerificationGasLimit.Uint64(), big.NewInt(f).Uint64(), nil
		}
		return retryEstimateGas(err, simOp.VerificationGasLimit.Int64(), in)
	}
	return simOp.VerificationGasLimit.Uint64(), simOp.CallGasLimit.Uint64(), nil
}
