package config

import (
	"math/big"

	mapset "github.com/deckarep/golang-set/v2"
)

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
	LyraChainID            = big.NewInt(957)
	LyraSepoliaChainID     = big.NewInt(902)
	Ancient8SepoliaChainID = big.NewInt(28122024)

	OpStackChains = mapset.NewSet(
		OptimismChainID.Uint64(),
		OptimismGoerliChainID.Uint64(),
		OptimismSepoliaChainID.Uint64(),
		BaseChainID.Uint64(),
		BaseGoerliChainID.Uint64(),
		BaseSepoliaChainID.Uint64(),
		LyraChainID.Uint64(),
		LyraSepoliaChainID.Uint64(),
		Ancient8SepoliaChainID.Uint64(),
	)

	ArbStackChains = mapset.NewSet(
		ArbitrumOneChainID.Uint64(),
		ArbitrumGoerliChainID.Uint64(),
		ArbitrumSepoliaChainID.Uint64(),
	)
)
