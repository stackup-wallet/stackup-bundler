package gas

import "math/big"

type GasEstimates struct {
	PreVerificationGas *big.Int `json:"preVerificationGas"`
	CallGasLimit       *big.Int `json:"callGasLimit"`
	VerificationGas    *big.Int `json:"verificationGas"`
}
