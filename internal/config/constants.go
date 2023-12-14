package config

import "math/big"

var (
	EthereumChainID        = big.NewInt(1)
	GoerliChainID          = big.NewInt(5)
	SepoliaChainID         = big.NewInt(11155111)
	ArbitrumOneChainID     = big.NewInt(42161)
	ArbitrumGoerliChainID  = big.NewInt(421613)
	ArbitrumSepoliaChainID = big.NewInt(421614)
	OptimismChainID        = big.NewInt(10)
	OptimismGoerliChainID  = big.NewInt(420)
	OptimismSepoliaChainID = big.NewInt(11155420)
	BaseChainID            = big.NewInt(8453)
	BaseGoerliChainID      = big.NewInt(84531)
	BaseSepoliaChainID     = big.NewInt(84532)
)
