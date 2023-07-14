package gasprice

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

// GetLegacyGasPriceFunc provides a general interface for retrieving the closest estimate for gas price to
// allow for timely execution of a transaction.
type GetLegacyGasPriceFunc = func() (*big.Int, error)

// NoopGetLegacyGasPriceFunc returns nil gas price and nil error.
func NoopGetLegacyGasPriceFunc() GetLegacyGasPriceFunc {
	return func() (*big.Int, error) {
		return nil, nil
	}
}

// GetLegacyGasPriceWithEthClient returns a GetLegacyGasPriceFunc using an eth client.
func GetLegacyGasPriceWithEthClient(eth *ethclient.Client) GetLegacyGasPriceFunc {
	return func() (*big.Int, error) {
		gp, err := eth.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, err
		}
		return gp, nil
	}
}
