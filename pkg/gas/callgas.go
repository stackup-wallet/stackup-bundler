package gas

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func CallGasEstimate(
	eth *ethclient.Client,
	from common.Address,
	op *userop.UserOperation,
) (uint64, error) {
	est, err := eth.EstimateGas(context.Background(), ethereum.CallMsg{
		From: from,
		To:   &op.Sender,
		Data: op.CallData,
	})
	if err != nil {
		return 0, errors.NewRPCError(errors.EXECUTION_REVERTED, err.Error(), err.Error())
	}

	return est, nil
}
