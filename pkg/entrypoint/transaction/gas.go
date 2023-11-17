package transaction

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// SuggestMeanGasTipCap suggests a Max Priority Fee for an EIP-1559 transaction to submit a batch of
// UserOperations to the EntryPoint. It returns the larger value between the suggested gas tip or the average
// maxPriorityFeePerGas weighted by max available gas per operation of the entire batch.
func SuggestMeanGasTipCap(tip *big.Int, batch []*userop.UserOperation) *big.Int {
	totalGasTip, totalGas := common.Big0, common.Big0
	for _, op := range batch {
		maxOpGas := op.GetMaxGasAvailable()
		totalGas = big.NewInt(0).Add(totalGas, maxOpGas)
		totalGasTip = big.NewInt(0).Add(
			totalGasTip, big.NewInt(0).Mul(maxOpGas, op.MaxPriorityFeePerGas),
		)
	}

	avg := common.Big0
	if totalGas.Cmp(common.Big0) != 0 {
		avg = big.NewInt(0).Div(totalGasTip, totalGas)
	}

	if avg.Cmp(tip) == 1 {
		return avg
	}
	return tip
}

// SuggestMeanGasFeeCap suggests a Max Fee for an EIP-1559 transaction to submit a batch of UserOperations to
// the EntryPoint. It returns the larger value between the recommended max fee or the average maxFeePerGas
// weighted by max available gas per operation of the entire batch.
func SuggestMeanGasFeeCap(basefee *big.Int, tip *big.Int, batch []*userop.UserOperation) *big.Int {
	mf := big.NewInt(0).Add(tip, big.NewInt(0).Mul(basefee, common.Big2))
	totalGasFee, totalGas := common.Big0, common.Big0
	for _, op := range batch {
		maxOpGas := op.GetMaxGasAvailable()
		totalGas = big.NewInt(0).Add(totalGas, maxOpGas)
		totalGasFee = big.NewInt(0).Add(
			totalGasFee, big.NewInt(0).Mul(maxOpGas, op.MaxFeePerGas),
		)
	}

	avg := common.Big0
	if totalGas.Cmp(common.Big0) != 0 {
		avg = big.NewInt(0).Div(totalGasFee, totalGas)
	}

	if avg.Cmp(mf) == 1 {
		return avg
	}
	return mf
}

// SuggestMeanGasPrice suggests a Gas Price for a legacy transaction to submit a batch of UserOperations to
// the EntryPoint. It returns the larger value between a given gas price or the average maxFeePerGas weighted
// by max available gas per operation of the entire batch.
func SuggestMeanGasPrice(gasPrice *big.Int, batch []*userop.UserOperation) *big.Int {
	totalGasFee, totalGas := common.Big0, common.Big0
	for _, op := range batch {
		maxOpGas := op.GetMaxGasAvailable()
		totalGas = big.NewInt(0).Add(totalGas, maxOpGas)
		totalGasFee = big.NewInt(0).Add(
			totalGasFee, big.NewInt(0).Mul(maxOpGas, op.MaxFeePerGas),
		)
	}

	avg := common.Big0
	if totalGas.Cmp(common.Big0) != 0 {
		avg = big.NewInt(0).Div(totalGasFee, totalGas)
	}

	if avg.Cmp(gasPrice) == 1 {
		return avg
	}
	return gasPrice
}
