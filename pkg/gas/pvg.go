package gas

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stackup-wallet/stackup-bundler/pkg/arbitrum/nodeinterface"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/methods"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

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
) CalcPreVerificationGasFunc {
	pk, _ := crypto.GenerateKey()
	dummy, _ := signer.New(hexutil.Encode(crypto.FromECDSA(pk))[2:])
	return func(op *userop.UserOperation, static *big.Int) (*big.Int, error) {
		// Pack handleOps method inputs
		ho, err := methods.HandleOpsMethod.Inputs.Pack(
			[]entrypoint.UserOperation{entrypoint.UserOperation(*op)},
			dummy.Address,
		)
		if err != nil {
			return nil, err
		}

		// Encode function data for gasEstimateL1Component
		create := false
		if op.Nonce.Cmp(common.Big0) == 0 {
			create = true
		}
		ge, err := nodeinterface.GasEstimateL1ComponentMethod.Inputs.Pack(
			nodeinterface.ERC4337GasHelperAddress,
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
