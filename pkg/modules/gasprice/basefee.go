package gasprice

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

// GetBaseFeeFunc provides a general interface for retrieving the closest estimate for basefee to allow for
// timely execution of a transaction.
type GetBaseFeeFunc = func() (*big.Int, error)

// NoopGetBaseFeeFunc returns nil basefee and nil error.
func NoopGetBaseFeeFunc() GetBaseFeeFunc {
	return func() (*big.Int, error) {
		return nil, nil
	}
}

// GetBaseFeeWithEthClient returns a GetBaseFeeFunc using an eth client.
func GetBaseFeeWithEthClient(eth *ethclient.Client) GetBaseFeeFunc {
	return func() (*big.Int, error) {
		head, err := eth.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		return head.BaseFee, nil
	}
}
