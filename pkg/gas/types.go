package gas

import "math/big"

// GasEstimates provides estimate values for all gas fields in a UserOperation.
type GasEstimates struct {
	PreVerificationGas   *big.Int `json:"preVerificationGas"`
	VerificationGasLimit *big.Int `json:"verificationGasLimit"`
	CallGasLimit         *big.Int `json:"callGasLimit"`

	// TODO: Deprecate in v0.7
	VerificationGas *big.Int `json:"verificationGas"`
}
