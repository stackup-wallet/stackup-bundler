package gasprice

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

// GetGasTipFunc provides a general interface for retrieving the closest estimate for gas tip to allow for
// timely execution of a transaction.
type GetGasTipFunc = func() (*big.Int, error)

// NoopGetGasTipFunc returns nil gas tip and nil error.
func NoopGetGasTipFunc() GetGasTipFunc {
	return func() (*big.Int, error) {
		return nil, nil
	}
}

// GetGasTipWithEthClient returns a GetGasTipFunc using an eth client.
func GetGasTipWithEthClient(eth *ethclient.Client) GetGasTipFunc {
	return func() (*big.Int, error) {
		gt, err := eth.SuggestGasTipCap(context.Background())
		if err != nil {
			return nil, err
		}
		return gt, nil
	}
}
