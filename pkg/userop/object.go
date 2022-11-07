package userop

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// UserOperation is the transaction object for ERC-4337 smart contract accounts.
type UserOperation struct {
	Sender               common.Address `json:"sender" mapstructure:"sender" validate:"required"`
	Nonce                *big.Int       `json:"nonce" mapstructure:"nonce" validate:"required"`
	InitCode             []byte         `json:"initCode"  mapstructure:"initCode" validate:"required"`
	CallData             []byte         `json:"callData" mapstructure:"callData" validate:"required"`
	CallGasLimit         *big.Int       `json:"callGasLimit" mapstructure:"callGasLimit" validate:"required"`
	VerificationGasLimit *big.Int       `json:"verificationGasLimit" mapstructure:"verificationGasLimit" validate:"required"`
	PreVerificationGas   *big.Int       `json:"preVerificationGas" mapstructure:"preVerificationGas" validate:"required"`
	MaxFeePerGas         *big.Int       `json:"maxFeePerGas" mapstructure:"maxFeePerGas" validate:"required"`
	MaxPriorityFeePerGas *big.Int       `json:"maxPriorityFeePerGas" mapstructure:"maxPriorityFeePerGas" validate:"required"`
	PaymasterAndData     []byte         `json:"paymasterAndData" mapstructure:"paymasterAndData" validate:"required"`
	Signature            []byte         `json:"signature" mapstructure:"signature" validate:"required"`
}

// Pack returns a standardized message of the op.
func (op *UserOperation) Pack() []byte {
	userOpType, _ := abi.NewType("tuple", "userOp", []abi.ArgumentMarshaling{
		{Name: "Sender", Type: "address"},
		{Name: "Nonce", Type: "uint256"},
		{Name: "InitCode", Type: "bytes"},
		{Name: "CallData", Type: "bytes"},
		{Name: "CallGasLimit", Type: "uint256"},
		{Name: "VerificationGasLimit", Type: "uint256"},
		{Name: "PreVerificationGas", Type: "uint256"},
		{Name: "MaxFeePerGas", Type: "uint256"},
		{Name: "MaxPriorityFeePerGas", Type: "uint256"},
		{Name: "PaymasterAndData", Type: "bytes"},
		{Name: "Signature", Type: "bytes"},
	})
	args := abi.Arguments{
		{Name: "UserOp", Type: userOpType},
	}
	packed, _ := args.Pack(&struct {
		Sender               common.Address
		Nonce                *big.Int
		InitCode             []byte
		CallData             []byte
		CallGasLimit         *big.Int
		VerificationGasLimit *big.Int
		PreVerificationGas   *big.Int
		MaxFeePerGas         *big.Int
		MaxPriorityFeePerGas *big.Int
		PaymasterAndData     []byte
		Signature            []byte
	}{
		op.Sender,
		op.Nonce,
		op.InitCode,
		op.CallData,
		op.CallGasLimit,
		op.VerificationGasLimit,
		op.PreVerificationGas,
		op.MaxFeePerGas,
		op.MaxPriorityFeePerGas,
		op.PaymasterAndData,
		[]byte{},
	})

	// Return with stripped leading word (total length) and trailing word (zero-length signature).
	enc := hexutil.Encode(packed)
	enc = "0x" + enc[66:len(enc)-64]
	return (hexutil.MustDecode(enc))
}

// GetRequestID returns the hash of op + entryPoint address + chainID.
func (op *UserOperation) GetRequestID(epAddr common.Address, chainID *big.Int) common.Hash {
	return crypto.Keccak256Hash(
		crypto.Keccak256(op.Pack()),
		common.LeftPadBytes(epAddr.Bytes(), 32),
		common.LeftPadBytes(chainID.Bytes(), 32),
	)
}

// MarshalJSON returns a JSON encoding of the UserOperation.
func (op *UserOperation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Sender               string `json:"sender"`
		Nonce                string `json:"nonce"`
		InitCode             string `json:"initCode"`
		CallData             string `json:"callData"`
		CallGasLimit         string `json:"callGasLimit"`
		VerificationGasLimit string `json:"verificationGasLimit"`
		PreVerificationGas   string `json:"preVerificationGas"`
		MaxFeePerGas         string `json:"maxFeePerGas"`
		MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
		PaymasterAndData     string `json:"paymasterAndData"`
		Signature            string `json:"signature"`
	}{
		Sender:               op.Sender.String(),
		Nonce:                op.Nonce.String(),
		InitCode:             hexutil.Encode(op.InitCode),
		CallData:             hexutil.Encode(op.CallData),
		CallGasLimit:         op.CallGasLimit.String(),
		VerificationGasLimit: op.CallGasLimit.String(),
		PreVerificationGas:   op.PreVerificationGas.String(),
		MaxFeePerGas:         op.MaxFeePerGas.String(),
		MaxPriorityFeePerGas: op.MaxPriorityFeePerGas.String(),
		PaymasterAndData:     hexutil.Encode(op.PaymasterAndData),
		Signature:            hexutil.Encode(op.Signature),
	})
}
