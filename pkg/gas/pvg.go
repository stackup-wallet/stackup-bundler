package gas

import (
	"bytes"
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/arbitrum/nodeinterface"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/methods"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/transaction"
	"github.com/stackup-wallet/stackup-bundler/pkg/optimism/gaspriceoracle"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// CalcPreVerificationGasFunc defines an interface for a function to calculate PVG given a userOp and a static
// value. The static input is the value derived from the default overheads.
type CalcPreVerificationGasFunc = func(op *userop.UserOperation, static *big.Int) (*big.Int, error)

func calcPVGFuncNoop() CalcPreVerificationGasFunc {
	return func(op *userop.UserOperation, static *big.Int) (*big.Int, error) {
		return nil, nil
	}
}

// CalcArbitrumPVGWithEthClient uses Arbitrum's NodeInterface precompile to get an estimate for
// preVerificationGas that takes into account the L1 gas component. see
// https://medium.com/offchainlabs/understanding-arbitrum-2-dimensional-fees-fd1d582596c9.
func CalcArbitrumPVGWithEthClient(
	rpc *rpc.Client,
	entryPoint common.Address,
) CalcPreVerificationGasFunc {
	pk, _ := crypto.GenerateKey()
	dummy, _ := signer.New(hexutil.Encode(crypto.FromECDSA(pk))[2:])
	return func(op *userop.UserOperation, static *big.Int) (*big.Int, error) {
		// Sanitize paymasterAndData.
		// TODO: Figure out why variability in this field is causing Arbitrum's precompile to return different
		// values.
		data, err := op.ToMap()
		if err != nil {
			return nil, err
		}
		data["paymasterAndData"] = hexutil.Encode(bytes.Repeat([]byte{1}, len(op.PaymasterAndData)))
		tmp, err := userop.New(data)
		if err != nil {
			return nil, err
		}

		// Pack handleOps method inputs
		ho, err := methods.HandleOpsMethod.Inputs.Pack(
			[]entrypoint.UserOperation{entrypoint.UserOperation(*tmp)},
			dummy.Address,
		)
		if err != nil {
			return nil, err
		}

		// Encode function data for gasEstimateL1Component
		create := false
		if tmp.Nonce.Cmp(common.Big0) == 0 {
			create = true
		}
		ge, err := nodeinterface.GasEstimateL1ComponentMethod.Inputs.Pack(
			entryPoint,
			create,
			append(methods.HandleOpsMethod.ID, ho...),
		)
		if err != nil {
			return nil, err
		}

		// Use eth_call to call the NodeInterface precompile
		req := map[string]any{
			"from": common.HexToAddress("0x"),
			"to":   nodeinterface.PrecompileAddress,
			"data": hexutil.Encode(append(nodeinterface.GasEstimateL1ComponentMethod.ID, ge...)),
		}
		var out any
		if err := rpc.Call(&out, "eth_call", &req, "latest"); err != nil {
			return nil, err
		}

		// Return static + GasEstimateForL1 as PVG
		gas, err := nodeinterface.DecodeGasEstimateL1ComponentOutput(out)
		if err != nil {
			return nil, err
		}
		return big.NewInt(0).Add(static, big.NewInt(int64(gas.GasEstimateForL1))), nil
	}
}

// CalcOptimismPVGWithEthClient uses Optimism's Gas Price Oracle precompile to get an estimate for
// preVerificationGas that takes into account the L1 gas component.
func CalcOptimismPVGWithEthClient(
	rpc *rpc.Client,
	chainID *big.Int,
	entryPoint common.Address,
) CalcPreVerificationGasFunc {
	pk, _ := crypto.GenerateKey()
	dummy, _ := signer.New(hexutil.Encode(crypto.FromECDSA(pk))[2:])
	return func(op *userop.UserOperation, static *big.Int) (*big.Int, error) {
		// Create Raw HandleOps Transaction
		eth := ethclient.NewClient(rpc)
		head, err := eth.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		tx, err := transaction.CreateRawHandleOps(&transaction.Opts{
			EOA:         dummy,
			Eth:         eth,
			ChainID:     chainID,
			EntryPoint:  entryPoint,
			Batch:       []*userop.UserOperation{op},
			Beneficiary: dummy.Address,
			BaseFee:     head.BaseFee,
			GasLimit:    math.MaxUint64,
		})
		if err != nil {
			return nil, err
		}

		// Encode function data for GetL1Fee
		data, err := hexutil.Decode(tx)
		if err != nil {
			return nil, err
		}
		ge, err := gaspriceoracle.GetL1FeeMethod.Inputs.Pack(data)
		if err != nil {
			return nil, err
		}

		// Use eth_call to call the Gas Price Oracle precompile
		req := map[string]any{
			"from": common.HexToAddress("0x"),
			"to":   gaspriceoracle.PrecompileAddress,
			"data": hexutil.Encode(append(gaspriceoracle.GetL1FeeMethod.ID, ge...)),
		}
		var out any
		if err := rpc.Call(&out, "eth_call", &req, "latest"); err != nil {
			return nil, err
		}

		// Get L1Fee and L2Price
		l1fee, err := gaspriceoracle.DecodeGetL1FeeMethodOutput(out)
		if err != nil {
			return nil, err
		}
		l2price := op.MaxFeePerGas
		l2priority := big.NewInt(0).Add(op.MaxPriorityFeePerGas, head.BaseFee)
		if l2priority.Cmp(l2price) == -1 {
			l2price = l2priority
		}

		// Return static + L1 buffer as PVG. L1 buffer is equal to L1Fee/L2Price.
		return big.NewInt(0).Add(static, big.NewInt(0).Div(l1fee, l2price)), nil
	}
}
