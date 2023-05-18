// Package gas implements helper functions for calculating gas parameters based on Ethereum protocol values.
package gas

import (
	"bytes"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Overhead provides helper methods for calculating gas limits based on pre-defined parameters.
type Overhead struct {
	fixed               float64
	perUserOp           float64
	perUserOpWord       float64
	zeroByte            float64
	nonZeroByte         float64
	minBundleSize       float64
	warmStorageRead     float64
	nonZeroValueCall    float64
	callOpcode          float64
	nonZeroValueStipend float64
	sanitizedPVG        *big.Int
	sanitizedVGL        *big.Int
	sanitizedCGL        *big.Int
	calcPVGFunc         CalcPreVerificationGasFunc
}

// NewDefaultOverhead returns an instance of Overhead using parameters defined by the Ethereum protocol.
func NewDefaultOverhead() *Overhead {
	return &Overhead{
		fixed:               21000,
		perUserOp:           18300,
		perUserOpWord:       4,
		zeroByte:            4,
		nonZeroByte:         16,
		minBundleSize:       1,
		warmStorageRead:     100,
		nonZeroValueCall:    9000,
		callOpcode:          700,
		nonZeroValueStipend: 2300,
		sanitizedPVG:        big.NewInt(100000),
		sanitizedVGL:        big.NewInt(1000000),
		sanitizedCGL:        big.NewInt(1000000),
		calcPVGFunc:         calcPVGFuncNoop(),
	}
}

// SetCalcPreVerificationGasFunc allows a custom function to be defined that can control how it calculates
// PVG. This is useful for networks that have different models for gas.
func (ov *Overhead) SetCalcPreVerificationGasFunc(fn CalcPreVerificationGasFunc) {
	ov.calcPVGFunc = fn
}

// CalcPreVerificationGas returns an expected gas cost for processing a UserOperation from a batch.
func (ov *Overhead) CalcPreVerificationGas(op *userop.UserOperation) (*big.Int, error) {
	// Sanitize fields to reduce as much variability due to length and zero bytes
	data, err := op.ToMap()
	if err != nil {
		return nil, err
	}
	data["preVerificationGas"] = hexutil.EncodeBig(ov.sanitizedPVG)
	data["verificationGasLimit"] = hexutil.EncodeBig(ov.sanitizedVGL)
	data["callGasLimit"] = hexutil.EncodeBig(ov.sanitizedCGL)
	data["signature"] = hexutil.Encode(bytes.Repeat([]byte{1}, len(op.Signature)))
	tmp, err := userop.New(data)
	if err != nil {
		return nil, err
	}

	// Calculate static value from pre-defined parameters
	packed := tmp.Pack()
	lengthInWord := float64(len(packed)+31) / 32
	callDataCost := float64(0)
	for _, b := range packed {
		if b == byte(0) {
			callDataCost += ov.zeroByte
		} else {
			callDataCost += ov.nonZeroByte
		}
	}
	pvg := callDataCost + (ov.fixed / ov.minBundleSize) + ov.perUserOp + (ov.perUserOpWord * lengthInWord)
	static := big.NewInt(int64(math.Round(pvg)))

	// Use value from CalcPreVerificationGasFunc if set, otherwise return the static value.
	g, err := ov.calcPVGFunc(tmp, static)
	if err != nil {
		return nil, err
	}
	if g != nil {
		return g, nil
	}
	return static, nil
}

// NonZeroValueCall returns an expected gas cost of using the CALL opcode in the context of EIP-4337.
// See https://github.com/wolflo/evm-opcodes/blob/main/gas.md#aa-1-call.
func (ov *Overhead) NonZeroValueCall() *big.Int {
	return big.NewInt(
		int64(ov.fixed + ov.warmStorageRead + ov.nonZeroValueCall + ov.callOpcode + ov.nonZeroValueStipend),
	)
}
