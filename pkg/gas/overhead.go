// Package gas implements helper functions for calculating gas parameters based on Ethereum protocol values.
package gas

import (
	"math"
	"math/big"

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
}

// CalcPreVerificationGas returns an expected gas cost for processing a UserOperation from a batch.
func (ov *Overhead) CalcPreVerificationGas(op *userop.UserOperation) *big.Int {
	packed := op.Pack()
	callDataCost := float64(0)

	for _, b := range packed {
		if b == byte(0) {
			callDataCost += ov.zeroByte
		} else {
			callDataCost += ov.nonZeroByte
		}
	}

	pvg := callDataCost + (ov.fixed / ov.minBundleSize) + ov.perUserOp + ov.perUserOpWord*float64(
		(len(packed)),
	)
	return big.NewInt(int64(math.Round(pvg)))
}

// NonZeroValueCall returns an expected gas cost of using the CALL opcode in the context of EIP-4337.
// See https://github.com/wolflo/evm-opcodes/blob/main/gas.md#aa-1-call.
func (ov *Overhead) NonZeroValueCall() *big.Int {
	return big.NewInt(
		int64(ov.warmStorageRead + ov.nonZeroValueCall + ov.callOpcode + ov.nonZeroValueStipend),
	)
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
	}
}
