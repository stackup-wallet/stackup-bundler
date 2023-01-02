package gas

import "math/big"

// GasEstimates provides estimate values for all gas fields in a UserOperation.
type GasEstimates struct {
	PreVerificationGas *big.Int `json:"preVerificationGas"`
	VerificationGas    *big.Int `json:"verificationGas"`
	CallGasLimit       *big.Int `json:"callGasLimit"`
}
