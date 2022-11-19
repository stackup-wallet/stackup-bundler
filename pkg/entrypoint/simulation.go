package entrypoint

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// SimulateValidationResults returns the results from a userOp simulation. For details see the inline docs
// for EntryPoint.sol#simulateValidation at https://github.com/eth-infinitism/account-abstraction.
type SimulateValidationResults struct {
	PreOpGas          *big.Int
	Prefund           *big.Int
	ActualAggregator  common.Address
	SigForUserOp      []byte
	SigForAggregation []byte
	OffChainSigInfo   []byte
}

// SimulateValidation makes a static call to Entrypoint.simulateValidation(userop, false) and returns the
// results without any state changes.
func SimulateValidation(eth *ethclient.Client, entryPoint common.Address, op *userop.UserOperation) (*SimulateValidationResults, error) {
	ep, err := NewEntrypoint(entryPoint, eth)
	if err != nil {
		return nil, err
	}

	var res []interface{}
	rawCaller := &EntrypointRaw{Contract: ep}
	err = rawCaller.Call(nil, &res, "simulateValidation", UserOperation(*op), false)
	if err != nil {
		revert, err := newFailedOpRevert(err)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(revert.Reason)
	}

	// TODO: Trace forbidden opcodes

	return &SimulateValidationResults{
		PreOpGas:          res[0].(*big.Int),
		Prefund:           res[1].(*big.Int),
		ActualAggregator:  res[2].(common.Address),
		SigForUserOp:      res[3].([]byte),
		SigForAggregation: res[4].([]byte),
		OffChainSigInfo:   res[5].([]byte),
	}, nil
}
