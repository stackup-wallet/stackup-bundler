package entrypoint

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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
func SimulateValidation(ep *Entrypoint, op UserOperation) (*SimulateValidationResults, error) {
	var res []interface{}
	rawCaller := &EntrypointRaw{Contract: ep}
	err := rawCaller.Call(nil, &res, "simulateValidation", op, false)
	if err != nil {
		return nil, err
	}

	return &SimulateValidationResults{
		PreOpGas:          res[0].(*big.Int),
		Prefund:           res[1].(*big.Int),
		ActualAggregator:  res[2].(common.Address),
		SigForUserOp:      res[3].([]byte),
		SigForAggregation: res[4].([]byte),
		OffChainSigInfo:   res[5].([]byte),
	}, nil
}
