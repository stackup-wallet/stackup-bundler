package gas

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type CalcPreVerificationGasFunc = func(op *userop.UserOperation) *big.Int

func calcPVGFuncNoop() CalcPreVerificationGasFunc {
	return func(op *userop.UserOperation) *big.Int {
		return nil
	}
}

func CalcArbitrumPVGWithEthClient(eth *ethclient.Client) CalcPreVerificationGasFunc {
	return func(op *userop.UserOperation) *big.Int {
		return big.NewInt(0)
	}
}
