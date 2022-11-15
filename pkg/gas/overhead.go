package gas

import (
	"math"
	"math/big"

	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Overhead provides helper methods for calculating gas limits based on pre-defined parameters.
type Overhead struct {
	Fixed         float64
	PerUserOp     float64
	PerUserOpWord float64
	ZeroByte      float64
	NonZeroByte   float64
	MinBundleSize float64
}

// CalcPreVerificationGas returns an expected gas cost for processing a UserOperation from a batch.
func (ov *Overhead) CalcPreVerificationGas(op *userop.UserOperation) *big.Int {
	packed := op.Pack()
	callDataCost := float64(0)

	for _, b := range packed {
		if b == byte(0) {
			callDataCost += ov.ZeroByte
		} else {
			callDataCost += ov.NonZeroByte
		}
	}

	pvg := callDataCost + (ov.Fixed / ov.MinBundleSize) + ov.PerUserOp + ov.PerUserOpWord*float64((len(packed)))
	return big.NewInt(int64(math.Round(pvg)))
}

// NewDefaultOverhead returns an instance of Overhead using parameters defined by the Ethereum protocol.
func NewDefaultOverhead() *Overhead {
	return &Overhead{
		Fixed:         21000,
		PerUserOp:     18300,
		PerUserOpWord: 4,
		ZeroByte:      4,
		NonZeroByte:   16,
		MinBundleSize: 1,
	}
}
