package execution

import (
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/utils"
	"github.com/stackup-wallet/stackup-bundler/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func TraceSimulateHandleOp(
	rpc *rpc.Client,
	entryPoint common.Address,
	op *userop.UserOperation,
	chainID *big.Int,
	customTracer string,
	target common.Address,
	data []byte,
) error {
	ep, err := entrypoint.NewEntrypoint(entryPoint, ethclient.NewClient(rpc))
	if err != nil {
		return err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(utils.DummyPk, chainID)
	if err != nil {
		return err
	}
	auth.GasLimit = math.MaxUint64
	auth.NoSend = true
	tx, err := ep.SimulateHandleOp(auth, entrypoint.UserOperation(*op), target, data)
	if err != nil {
		return err
	}

	var res tracer.BundlerErrorReturn
	req := utils.TraceCallReq{
		From: common.HexToAddress("0x"),
		To:   entryPoint,
		Data: tx.Data(),
	}
	opts := utils.TraceCallOpts{
		Tracer: customTracer,
	}
	if err := rpc.CallContext(context.Background(), &res, "debug_traceCall", &req, "latest", &opts); err != nil {
		return err
	}

	if len(res.Reverts) != 0 {
		data, err := hexutil.Decode(res.Reverts[len(res.Reverts)-1])
		if err != nil {
			return err
		}

		reason, err := errors.DecodeRevert(data)
		if err != nil {
			return err
		}
		return errors.NewRPCError(errors.EXECUTION_REVERTED, reason, reason)
	}
	return nil
}
