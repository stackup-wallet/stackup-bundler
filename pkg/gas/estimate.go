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
	var simErr error
	for l <= r {
		m := (l + r) / 2

		data["verificationGasLimit"] = hexutil.EncodeBig(big.NewInt(int64(m)))
		simOp, err := userop.New(data)
		if err != nil {
			return 0, 0, err
		}
		out, err := execution.TraceSimulateHandleOp(&execution.TraceInput{
			Rpc:        in.Rpc,
			EntryPoint: in.EntryPoint,
			Op:         simOp,
			ChainID:    in.ChainID,
		})
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
				return 0, 0, err
			}
		}

		// Optimal VGL found.
		data["verificationGasLimit"] = hexutil.EncodeBig(
			big.NewInt(0).Sub(out.Result.PreOpGas, in.Op.PreVerificationGas),
		)
		break
	}
	if simErr != nil && !isExecutionOOG(simErr) {
		return 0, 0, simErr
	}

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
		Rpc:        in.Rpc,
		EntryPoint: in.EntryPoint,
		Op:         simOp,
		ChainID:    in.ChainID,
	})
	if err != nil {
		return 0, 0, err
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
		if isExecutionOOG(err) {
			l := cgl.Int64()
			r := l * 2
			count := 0
			var simErr error
			for r-l >= fallBackBinarySearchCutoff {
				count++
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
				if err != nil && (isExecutionOOG(err) || isExecutionReverted(err)) {
					// CGL too low, go higher.
					l = m + 1
					continue
				} else if err == nil || isPrefundNotPaid(err) {
					// CGL too high, go lower.
					r = m - 1
					continue
				} else {
					// Unexpected error.
					return 0, 0, err
				}
			}
			if simErr != nil && !(isExecutionOOG(err) || isExecutionReverted(err)) {
				return 0, 0, simErr
			}
			return simOp.VerificationGasLimit.Uint64(), big.NewInt(r).Uint64(), nil
		}
		return 0, 0, err
	}
	return simOp.VerificationGasLimit.Uint64(), simOp.CallGasLimit.Uint64(), nil
}
