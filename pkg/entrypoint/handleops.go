package entrypoint

import (
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func ToAbiType(batch []*userop.UserOperation) []UserOperation {
	ops := []UserOperation{}
	for _, op := range batch {
		ops = append(ops, UserOperation(*op))
	}

	return ops
}

// EstimateHandleOpsGas returns a gas estimate required to call handleOps() with a given batch. A failed call
// will return the cause of the revert.
func EstimateHandleOpsGas(
	eoa *signer.EOA,
	eth *ethclient.Client,
	chainID *big.Int,
	entryPoint common.Address,
	batch []*userop.UserOperation,
	beneficiary common.Address,
) (gas uint64, revert *FailedOpRevert, err error) {
	ep, err := NewEntrypoint(entryPoint, eth)
	if err != nil {
		return 0, nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(eoa.PrivateKey, chainID)
	if err != nil {
		return 0, nil, err
	}
	auth.GasLimit = math.MaxUint64
	auth.NoSend = true

	tx, err := ep.HandleOps(auth, ToAbiType(batch), beneficiary)
	if err != nil {
		return 0, nil, err
	}

	est, err := eth.EstimateGas(context.Background(), ethereum.CallMsg{
		From:       eoa.Address,
		To:         tx.To(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		GasFeeCap:  tx.GasFeeCap(),
		GasTipCap:  tx.GasTipCap(),
		Value:      tx.Value(),
		Data:       tx.Data(),
		AccessList: tx.AccessList(),
	})
	if err != nil {
		revert, err := newFailedOpRevert(err)
		if err != nil {
			return 0, nil, err
		}
		return 0, revert, nil
	}

	return est, nil, nil
}

func HandleOps(
	eoa *signer.EOA,
	eth *ethclient.Client,
	chainID *big.Int,
	entryPoint common.Address,
	batch []*userop.UserOperation,
	beneficiary common.Address,
	gas uint64,
	tip *big.Int,
) (revert *FailedOpRevert, err error) {
	ep, err := NewEntrypoint(entryPoint, eth)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(eoa.PrivateKey, chainID)
	if err != nil {
		return nil, err
	}
	auth.GasLimit = gas
	auth.GasTipCap = tip

	_, err = ep.HandleOps(auth, ToAbiType(batch), beneficiary)
	if err != nil {
		revert, err := newFailedOpRevert(err)
		if err != nil {
			return nil, err
		}
		return revert, nil
	}

	return nil, nil
}
