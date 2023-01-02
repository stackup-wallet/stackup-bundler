package gas

import "math/big"

type GasEstimates struct {
	PreVerificationGas *big.Int `json:"preVerificationGas"`
	VerificationGas    *big.Int `json:"verificationGas"`
	CallGasLimit       *big.Int `json:"callGasLimit"`
}
