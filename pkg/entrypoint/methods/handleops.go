package methods

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

var (
	HandleOpsMethod = abi.NewMethod(
		"handleOps",
		"handleOps",
		abi.Function,
		"",
		false,
		false,
		abi.Arguments{
			{Name: "ops", Type: userop.UserOpArr},
			{Name: "beneficiary", Type: address},
		},
		nil,
	)
	HandleOpsSelector = hexutil.Encode(HandleOpsMethod.ID)
)
