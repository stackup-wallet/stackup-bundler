package transaction

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// SuggestMeanGasTipCap suggests a Max Priority Fee for an EIP-1559 transaction to submit a batch of
// UserOperations to the EntryPoint. It returns the larger value between the suggested gas tip or the average
// maxPriorityFeePerGas of the entire batch.
func SuggestMeanGasTipCap(tip *big.Int, batch []*userop.UserOperation) *big.Int {
	sum := big.NewInt(0)
	for _, op := range batch {
		sum = big.NewInt(0).Add(sum, op.MaxPriorityFeePerGas)
	}
	avg := big.NewInt(0).Div(sum, big.NewInt(int64(len(batch))))

	if avg.Cmp(tip) == 1 {
		return avg
	}
	return tip
}

// SuggestMeanGasFeeCap suggests a Max Fee for an EIP-1559 transaction to submit a batch of UserOperations to
// the EntryPoint. It returns the larger value between the recommended max fee or the average maxFeePerGas of
// the entire batch.
func SuggestMeanGasFeeCap(basefee *big.Int, tip *big.Int, batch []*userop.UserOperation) *big.Int {
	mf := big.NewInt(0).Add(tip, big.NewInt(0).Mul(basefee, common.Big2))
	sum := big.NewInt(0)
	for _, op := range batch {
		sum = big.NewInt(0).Add(sum, op.MaxFeePerGas)
	}
	avg := big.NewInt(0).Div(sum, big.NewInt(int64(len(batch))))

	if avg.Cmp(mf) == 1 {
		return avg
	}
	return mf
}

// SuggestMeanGasPrice suggests a Gas Price for a legacy transaction to submit a batch of UserOperations to
// the EntryPoint. It returns the larger value between a given gas price or the average maxFeePerGas of the
// entire batch.
func SuggestMeanGasPrice(gasPrice *big.Int, batch []*userop.UserOperation) *big.Int {
	sum := big.NewInt(0)
	for _, op := range batch {
		sum = big.NewInt(0).Add(sum, op.MaxFeePerGas)
	}
	avg := big.NewInt(0).Div(sum, big.NewInt(int64(len(batch))))

	if avg.Cmp(gasPrice) == 1 {
		return avg
	}
	return gasPrice
}
