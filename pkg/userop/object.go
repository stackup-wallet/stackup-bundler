// Package userop provides the base transaction object used throughout the stackup-bundler.
package userop

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// UserOpType is the ABI type of a UserOperation.
	UserOpType, _ = abi.NewType("tuple", "userOp", []abi.ArgumentMarshaling{
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
)

func getAbiArgs() abi.Arguments {
	return abi.Arguments{
		{Name: "UserOp", Type: UserOpType},
	}
}

// UserOperation represents an EIP-4337 style transaction for a smart contract account.
type UserOperation struct {
	Sender               common.Address `json:"sender"               mapstructure:"sender"               validate:"required"`
	Nonce                *big.Int       `json:"nonce"                mapstructure:"nonce"                validate:"required"`
	InitCode             []byte         `json:"initCode"             mapstructure:"initCode"             validate:"required"`
	CallData             []byte         `json:"callData"             mapstructure:"callData"             validate:"required"`
	CallGasLimit         *big.Int       `json:"callGasLimit"         mapstructure:"callGasLimit"         validate:"required"`
	VerificationGasLimit *big.Int       `json:"verificationGasLimit" mapstructure:"verificationGasLimit" validate:"required"`
	PreVerificationGas   *big.Int       `json:"preVerificationGas"   mapstructure:"preVerificationGas"   validate:"required"`
	MaxFeePerGas         *big.Int       `json:"maxFeePerGas"         mapstructure:"maxFeePerGas"         validate:"required"`
	MaxPriorityFeePerGas *big.Int       `json:"maxPriorityFeePerGas" mapstructure:"maxPriorityFeePerGas" validate:"required"`
	PaymasterAndData     []byte         `json:"paymasterAndData"     mapstructure:"paymasterAndData"     validate:"required"`
	Signature            []byte         `json:"signature"            mapstructure:"signature"            validate:"required"`
}

// GetPaymaster returns the address portion of PaymasterAndData if applicable. Otherwise it returns the zero
// address.
func (op *UserOperation) GetPaymaster() common.Address {
	if len(op.PaymasterAndData) < common.AddressLength {
		return common.HexToAddress("0x")
	}

	return common.BytesToAddress(op.PaymasterAndData[:common.AddressLength])
}

// GetFactory returns the address portion of InitCode if applicable. Otherwise it returns the zero address.
func (op *UserOperation) GetFactory() common.Address {
	if len(op.InitCode) < common.AddressLength {
		return common.HexToAddress("0x")
	}

	return common.BytesToAddress(op.InitCode[:common.AddressLength])
}

// GetMaxPrefund returns the max amount of wei required to pay for gas fees by either the sender or
// paymaster.
func (op *UserOperation) GetMaxPrefund() *big.Int {
	mul := big.NewInt(1)
	paymaster := op.GetPaymaster()
	if paymaster != common.HexToAddress("0x") {
		mul = big.NewInt(3)
	}

	requiredGas := big.NewInt(0).Add(
		big.NewInt(0).Mul(op.VerificationGasLimit, mul),
		big.NewInt(0).Add(op.PreVerificationGas, op.CallGasLimit),
	)
	return big.NewInt(0).Mul(requiredGas, op.MaxFeePerGas)
}

// Pack returns a standard message of the userOp. This cannot be used to generate a userOpHash.
func (op *UserOperation) Pack() []byte {
	args := getAbiArgs()
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
		op.Signature,
	})

	enc := hexutil.Encode(packed)
	enc = "0x" + enc[66:]
	return (hexutil.MustDecode(enc))
}

// PackForSignature returns a minimal message of the userOp. This can be used to generate a userOpHash.
func (op *UserOperation) PackForSignature() []byte {
	args := getAbiArgs()
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

// GetUserOpHash returns the hash of the userOp + entryPoint address + chainID.
func (op *UserOperation) GetUserOpHash(entryPoint common.Address, chainID *big.Int) common.Hash {
	return crypto.Keccak256Hash(
		crypto.Keccak256(op.PackForSignature()),
		common.LeftPadBytes(entryPoint.Bytes(), 32),
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
		Nonce:                hexutil.EncodeBig(op.Nonce),
		InitCode:             hexutil.Encode(op.InitCode),
		CallData:             hexutil.Encode(op.CallData),
		CallGasLimit:         hexutil.EncodeBig(op.CallGasLimit),
		VerificationGasLimit: hexutil.EncodeBig(op.VerificationGasLimit),
		PreVerificationGas:   hexutil.EncodeBig(op.PreVerificationGas),
		MaxFeePerGas:         hexutil.EncodeBig(op.MaxFeePerGas),
		MaxPriorityFeePerGas: hexutil.EncodeBig(op.MaxPriorityFeePerGas),
		PaymasterAndData:     hexutil.Encode(op.PaymasterAndData),
		Signature:            hexutil.Encode(op.Signature),
	})
}
