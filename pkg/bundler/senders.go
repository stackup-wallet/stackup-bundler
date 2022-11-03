package bundler

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func getSenders(batch []*userop.UserOperation) []common.Address {
	s := []common.Address{}
	for _, op := range batch {
		s = append(s, op.Sender)
	}

	return s
}
