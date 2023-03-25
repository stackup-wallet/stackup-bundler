package guardian

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"math/big"
)

type Guardian struct {
	Eoa                      *signer.EOA
	ContractAddress          *common.Address
	PrivateAccountTransactor *PrivateRecoveryAccountTransactor
}

type RecoverRequest struct {
	NewOwner common.Address `json:"new_owner"`
	A        [2]*big.Int    `json:"a"`
	B        [2][2]*big.Int `json:"b"`
	C        [2]*big.Int    `json:"c"`
	Input    [3]*big.Int    `json:"input"`
}
