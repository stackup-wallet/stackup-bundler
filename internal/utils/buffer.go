package utils

import "math/big"

func AddBuffer(amt *big.Int, factor int64) *big.Int {
	if amt == nil {
		return nil
	}

	a := big.NewInt(0).Mul(amt, big.NewInt(1000+(factor*10)))
	return big.NewInt(0).Div(a, big.NewInt(1000))
}
