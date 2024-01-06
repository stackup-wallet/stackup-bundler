package checks

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetCodeFunc provides a general interface for retrieving the bytecode for a given address.
type GetCodeFunc = func(addr common.Address) ([]byte, error)

// getCodeWithEthClient returns a GetCodeFunc that uses an eth client to call eth_getCode.
func getCodeWithEthClient(eth *ethclient.Client) GetCodeFunc {
	return func(addr common.Address) ([]byte, error) {
		return eth.CodeAt(context.Background(), addr, nil)
	}
}
