package stake

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint"
)

// GetStakeFunc provides a general interface for retrieving the EntryPoint stake for a given address.
type GetStakeFunc = func(entryPoint, entity common.Address) (entrypoint.IStakeManagerDepositInfo, error)

func GetStakeFuncNoop(ver string) GetStakeFunc {
	return func(entryPoint, entity common.Address) (entrypoint.IStakeManagerDepositInfo, error) {
		return entrypoint.NewStakeManagerByVersion(ver), nil
	}
}

// GetStakeWithEthClient returns a GetStakeFunc that relies on an eth client to get stake info from the
// EntryPoint.
func GetStakeWithEthClient(eth *ethclient.Client) GetStakeFunc {
	return func(entryPoint, addr common.Address) (entrypoint.IStakeManagerDepositInfo, error) {
		if addr == common.HexToAddress("0x") {
			return nil, nil
		}

		ep, err := entrypoint.NewEntrypoint(entryPoint, eth)
		if err != nil {
			return nil, err
		}

		dep, err := ep.GetDepositInfo(nil, addr)
		if err != nil {
			return nil, err
		}

		return &dep, nil
	}
}
