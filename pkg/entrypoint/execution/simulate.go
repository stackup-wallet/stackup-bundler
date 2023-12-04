package execution

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/reverts"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/utils"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/state"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type SimulateInput struct {
	Rpc        *rpc.Client
	EntryPoint common.Address
	Op         *userop.UserOperation
	Sos        state.OverrideSet
	ChainID    *big.Int

	// Optional params for simulateHandleOps
	Target common.Address
	Data   []byte
}

func SimulateHandleOp(in *SimulateInput) (*reverts.ExecutionResultRevert, error) {
	ep, err := entrypoint.NewEntrypoint(in.EntryPoint, ethclient.NewClient(in.Rpc))
	if err != nil {
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(utils.DummyPk, in.ChainID)
	if err != nil {
		return nil, err
	}
	auth.GasLimit = math.MaxUint64
	auth.NoSend = true
	tx, err := ep.SimulateHandleOp(auth, entrypoint.UserOperation(*in.Op), in.Target, in.Data)
	if err != nil {
		return nil, err
	}

	req := utils.EthCallReq{
		From: common.HexToAddress("0x"),
		To:   in.EntryPoint,
		Data: tx.Data(),
	}
	err = in.Rpc.CallContext(context.Background(), nil, "eth_call", &req, "latest", in.Sos)

	sim, simErr := reverts.NewExecutionResult(err)
	if simErr != nil {
		fo, foErr := reverts.NewFailedOp(err)
		if foErr != nil {
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("%s, %s", simErr, foErr)
		}
		return nil, errors.NewRPCError(errors.REJECTED_BY_EP_OR_ACCOUNT, fo.Reason, fo)
	}

	return sim, nil
}
