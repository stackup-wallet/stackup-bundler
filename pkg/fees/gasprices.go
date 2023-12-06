package fees

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GasPrices contains recommended gas fees for a UserOperation to be included in a timely manner.
type GasPrices struct {
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
}

// NewGasPrices returns an instance of GasPrices with the latest suggested fees derived from an Eth Client.
func NewGasPrices(eth *ethclient.Client) (*GasPrices, error) {
	gp := GasPrices{}
	if head, err := eth.HeaderByNumber(context.Background(), nil); err != nil {
		return nil, err
	} else if head.BaseFee != nil {
		tip, err := eth.SuggestGasTipCap(context.Background())
		if err != nil {
			return nil, err
		}
		gp.MaxFeePerGas = big.NewInt(0).Add(tip, big.NewInt(0).Mul(head.BaseFee, common.Big2))
		gp.MaxPriorityFeePerGas = tip
	} else {
		sgp, err := eth.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, err
		}
		gp.MaxFeePerGas = sgp
		gp.MaxPriorityFeePerGas = sgp
	}

	return &gp, nil
}
