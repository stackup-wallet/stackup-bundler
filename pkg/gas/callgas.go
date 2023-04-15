package gas

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/execution"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func CallGasEstimate(
	eth *ethclient.Client,
	from common.Address,
	op *userop.UserOperation,
) (uint64, error) {
	data, err := op.ToMap()
	if err != nil {
		return 0, err
	}

	// Set MaxPriorityFeePerGas = MaxFeePerGas to simplify callGasLimit calculation from simulation paid
	// value.
	data["maxPriorityFeePerGas"] = hexutil.EncodeBig(op.MaxFeePerGas)
	simOp, err := userop.New(data)
	if err != nil {
		return 0, err
	}

	sim, err := execution.SimulateHandleOp(eth, from, simOp, common.Address{}, []byte{})
	if err != nil {
		return 0, err
	}
	cgl := big.NewInt(0).Div(sim.Paid, op.MaxFeePerGas)

	return cgl.Uint64(), nil
}
