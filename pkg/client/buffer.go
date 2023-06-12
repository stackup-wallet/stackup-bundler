package client

import "math/big"

func addBuffer(amt *big.Int, factor int64) *big.Int {
	a := big.NewInt(0).Mul(amt, big.NewInt(1000+(factor*10)))
	return big.NewInt(0).Div(a, big.NewInt(1000))
}
