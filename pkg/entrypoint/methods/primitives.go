package methods

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	bytes32, _ = abi.NewType("bytes32", "", nil)
	uint256, _ = abi.NewType("uint256", "", nil)
	bytes, _   = abi.NewType("bytes", "", nil)
	address, _ = abi.NewType("address", "", nil)
)
