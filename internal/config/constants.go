package config

import "math/big"

var (
	EthereumChainID       = big.NewInt(1)
	GoerliChainID         = big.NewInt(5)
	ArbitrumOneChainID    = big.NewInt(42161)
	ArbitrumGoerliChainID = big.NewInt(421613)
	OptimismChainID       = big.NewInt(10)
	OptimismGoerliChainID = big.NewInt(420)
	BaseChainID           = big.NewInt(8453)
	BaseGoerliChainID     = big.NewInt(84531)
)
