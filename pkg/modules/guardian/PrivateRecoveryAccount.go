// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package guardian

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// UserOperation is an auto generated low-level Go binding around an user-defined struct.
type UserOperation struct {
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
}

// PrivateRecoveryAccountMetaData contains all meta data concerning the PrivateRecoveryAccount contract.
var PrivateRecoveryAccountMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIEntryPoint\",\"name\":\"anEntryPoint\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIEntryPoint\",\"name\":\"entryPoint\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"PrivateRecoveryAccountInitialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"addDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"entryPoint\",\"outputs\":[{\"internalType\":\"contractIEntryPoint\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"func\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"dest\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"func\",\"type\":\"bytes[]\"}],\"name\":\"executeBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGuardians\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"anOwner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"guardians\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"vote_threshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"updateGuardianVerifierAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"socialRecoveryVerifierAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"poseidonContractAddress\",\"type\":\"address\"}],\"name\":\"initilizeGuardians\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"},{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[3]\",\"name\":\"input\",\"type\":\"uint256[3]\"}],\"name\":\"recover\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"valid\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"update\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[6]\",\"name\":\"input\",\"type\":\"uint256[6]\"}],\"name\":\"updateGuardian\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"verificationGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxPriorityFeePerGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"missingAccountFunds\",\"type\":\"uint256\"}],\"name\":\"validateUserOp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"withdrawAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawDepositTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// PrivateRecoveryAccountABI is the input ABI used to generate the binding from.
// Deprecated: Use PrivateRecoveryAccountMetaData.ABI instead.
var PrivateRecoveryAccountABI = PrivateRecoveryAccountMetaData.ABI

// PrivateRecoveryAccount is an auto generated Go binding around an Ethereum contract.
type PrivateRecoveryAccount struct {
	PrivateRecoveryAccountCaller     // Read-only binding to the contract
	PrivateRecoveryAccountTransactor // Write-only binding to the contract
	PrivateRecoveryAccountFilterer   // Log filterer for contract events
}

// PrivateRecoveryAccountCaller is an auto generated read-only Go binding around an Ethereum contract.
type PrivateRecoveryAccountCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrivateRecoveryAccountTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PrivateRecoveryAccountTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrivateRecoveryAccountFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PrivateRecoveryAccountFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrivateRecoveryAccountSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PrivateRecoveryAccountSession struct {
	Contract     *PrivateRecoveryAccount // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PrivateRecoveryAccountCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PrivateRecoveryAccountCallerSession struct {
	Contract *PrivateRecoveryAccountCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// PrivateRecoveryAccountTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PrivateRecoveryAccountTransactorSession struct {
	Contract     *PrivateRecoveryAccountTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// PrivateRecoveryAccountRaw is an auto generated low-level Go binding around an Ethereum contract.
type PrivateRecoveryAccountRaw struct {
	Contract *PrivateRecoveryAccount // Generic contract binding to access the raw methods on
}

// PrivateRecoveryAccountCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PrivateRecoveryAccountCallerRaw struct {
	Contract *PrivateRecoveryAccountCaller // Generic read-only contract binding to access the raw methods on
}

// PrivateRecoveryAccountTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PrivateRecoveryAccountTransactorRaw struct {
	Contract *PrivateRecoveryAccountTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPrivateRecoveryAccount creates a new instance of PrivateRecoveryAccount, bound to a specific deployed contract.
func NewPrivateRecoveryAccount(address common.Address, backend bind.ContractBackend) (*PrivateRecoveryAccount, error) {
	contract, err := bindPrivateRecoveryAccount(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PrivateRecoveryAccount{PrivateRecoveryAccountCaller: PrivateRecoveryAccountCaller{contract: contract}, PrivateRecoveryAccountTransactor: PrivateRecoveryAccountTransactor{contract: contract}, PrivateRecoveryAccountFilterer: PrivateRecoveryAccountFilterer{contract: contract}}, nil
}

// NewPrivateRecoveryAccountCaller creates a new read-only instance of PrivateRecoveryAccount, bound to a specific deployed contract.
func NewPrivateRecoveryAccountCaller(address common.Address, caller bind.ContractCaller) (*PrivateRecoveryAccountCaller, error) {
	contract, err := bindPrivateRecoveryAccount(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PrivateRecoveryAccountCaller{contract: contract}, nil
}

// NewPrivateRecoveryAccountTransactor creates a new write-only instance of PrivateRecoveryAccount, bound to a specific deployed contract.
func NewPrivateRecoveryAccountTransactor(address common.Address, transactor bind.ContractTransactor) (*PrivateRecoveryAccountTransactor, error) {
	contract, err := bindPrivateRecoveryAccount(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PrivateRecoveryAccountTransactor{contract: contract}, nil
}

// NewPrivateRecoveryAccountFilterer creates a new log filterer instance of PrivateRecoveryAccount, bound to a specific deployed contract.
func NewPrivateRecoveryAccountFilterer(address common.Address, filterer bind.ContractFilterer) (*PrivateRecoveryAccountFilterer, error) {
	contract, err := bindPrivateRecoveryAccount(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PrivateRecoveryAccountFilterer{contract: contract}, nil
}

// bindPrivateRecoveryAccount binds a generic wrapper to an already deployed contract.
func bindPrivateRecoveryAccount(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PrivateRecoveryAccountMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PrivateRecoveryAccount *PrivateRecoveryAccountRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PrivateRecoveryAccount.Contract.PrivateRecoveryAccountCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PrivateRecoveryAccount *PrivateRecoveryAccountRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.PrivateRecoveryAccountTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PrivateRecoveryAccount *PrivateRecoveryAccountRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.PrivateRecoveryAccountTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PrivateRecoveryAccount.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.contract.Transact(opts, method, params...)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCaller) EntryPoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PrivateRecoveryAccount.contract.Call(opts, &out, "entryPoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) EntryPoint() (common.Address, error) {
	return _PrivateRecoveryAccount.Contract.EntryPoint(&_PrivateRecoveryAccount.CallOpts)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCallerSession) EntryPoint() (common.Address, error) {
	return _PrivateRecoveryAccount.Contract.EntryPoint(&_PrivateRecoveryAccount.CallOpts)
}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCaller) GetDeposit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PrivateRecoveryAccount.contract.Call(opts, &out, "getDeposit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) GetDeposit() (*big.Int, error) {
	return _PrivateRecoveryAccount.Contract.GetDeposit(&_PrivateRecoveryAccount.CallOpts)
}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCallerSession) GetDeposit() (*big.Int, error) {
	return _PrivateRecoveryAccount.Contract.GetDeposit(&_PrivateRecoveryAccount.CallOpts)
}

// GetGuardians is a free data retrieval call binding the contract method 0x0665f04b.
//
// Solidity: function getGuardians() view returns(uint256[])
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCaller) GetGuardians(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _PrivateRecoveryAccount.contract.Call(opts, &out, "getGuardians")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetGuardians is a free data retrieval call binding the contract method 0x0665f04b.
//
// Solidity: function getGuardians() view returns(uint256[])
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) GetGuardians() ([]*big.Int, error) {
	return _PrivateRecoveryAccount.Contract.GetGuardians(&_PrivateRecoveryAccount.CallOpts)
}

// GetGuardians is a free data retrieval call binding the contract method 0x0665f04b.
//
// Solidity: function getGuardians() view returns(uint256[])
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCallerSession) GetGuardians() ([]*big.Int, error) {
	return _PrivateRecoveryAccount.Contract.GetGuardians(&_PrivateRecoveryAccount.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCaller) Nonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PrivateRecoveryAccount.contract.Call(opts, &out, "nonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) Nonce() (*big.Int, error) {
	return _PrivateRecoveryAccount.Contract.Nonce(&_PrivateRecoveryAccount.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCallerSession) Nonce() (*big.Int, error) {
	return _PrivateRecoveryAccount.Contract.Nonce(&_PrivateRecoveryAccount.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PrivateRecoveryAccount.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) Owner() (common.Address, error) {
	return _PrivateRecoveryAccount.Contract.Owner(&_PrivateRecoveryAccount.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCallerSession) Owner() (common.Address, error) {
	return _PrivateRecoveryAccount.Contract.Owner(&_PrivateRecoveryAccount.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PrivateRecoveryAccount.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) PendingOwner() (common.Address, error) {
	return _PrivateRecoveryAccount.Contract.PendingOwner(&_PrivateRecoveryAccount.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCallerSession) PendingOwner() (common.Address, error) {
	return _PrivateRecoveryAccount.Contract.PendingOwner(&_PrivateRecoveryAccount.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _PrivateRecoveryAccount.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) ProxiableUUID() ([32]byte, error) {
	return _PrivateRecoveryAccount.Contract.ProxiableUUID(&_PrivateRecoveryAccount.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountCallerSession) ProxiableUUID() ([32]byte, error) {
	return _PrivateRecoveryAccount.Contract.ProxiableUUID(&_PrivateRecoveryAccount.CallOpts)
}

// AddDeposit is a paid mutator transaction binding the contract method 0x4a58db19.
//
// Solidity: function addDeposit() payable returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactor) AddDeposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.contract.Transact(opts, "addDeposit")
}

// AddDeposit is a paid mutator transaction binding the contract method 0x4a58db19.
//
// Solidity: function addDeposit() payable returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) AddDeposit() (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.AddDeposit(&_PrivateRecoveryAccount.TransactOpts)
}

// AddDeposit is a paid mutator transaction binding the contract method 0x4a58db19.
//
// Solidity: function addDeposit() payable returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorSession) AddDeposit() (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.AddDeposit(&_PrivateRecoveryAccount.TransactOpts)
}

// Execute is a paid mutator transaction binding the contract method 0xb61d27f6.
//
// Solidity: function execute(address dest, uint256 value, bytes func) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactor) Execute(opts *bind.TransactOpts, dest common.Address, value *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.contract.Transact(opts, "execute", dest, value, arg2)
}

// Execute is a paid mutator transaction binding the contract method 0xb61d27f6.
//
// Solidity: function execute(address dest, uint256 value, bytes func) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) Execute(dest common.Address, value *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.Execute(&_PrivateRecoveryAccount.TransactOpts, dest, value, arg2)
}

// Execute is a paid mutator transaction binding the contract method 0xb61d27f6.
//
// Solidity: function execute(address dest, uint256 value, bytes func) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorSession) Execute(dest common.Address, value *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.Execute(&_PrivateRecoveryAccount.TransactOpts, dest, value, arg2)
}

// ExecuteBatch is a paid mutator transaction binding the contract method 0x18dfb3c7.
//
// Solidity: function executeBatch(address[] dest, bytes[] func) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactor) ExecuteBatch(opts *bind.TransactOpts, dest []common.Address, arg1 [][]byte) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.contract.Transact(opts, "executeBatch", dest, arg1)
}

// ExecuteBatch is a paid mutator transaction binding the contract method 0x18dfb3c7.
//
// Solidity: function executeBatch(address[] dest, bytes[] func) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) ExecuteBatch(dest []common.Address, arg1 [][]byte) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.ExecuteBatch(&_PrivateRecoveryAccount.TransactOpts, dest, arg1)
}

// ExecuteBatch is a paid mutator transaction binding the contract method 0x18dfb3c7.
//
// Solidity: function executeBatch(address[] dest, bytes[] func) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorSession) ExecuteBatch(dest []common.Address, arg1 [][]byte) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.ExecuteBatch(&_PrivateRecoveryAccount.TransactOpts, dest, arg1)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address anOwner) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactor) Initialize(opts *bind.TransactOpts, anOwner common.Address) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.contract.Transact(opts, "initialize", anOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address anOwner) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) Initialize(anOwner common.Address) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.Initialize(&_PrivateRecoveryAccount.TransactOpts, anOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address anOwner) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorSession) Initialize(anOwner common.Address) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.Initialize(&_PrivateRecoveryAccount.TransactOpts, anOwner)
}

// InitilizeGuardians is a paid mutator transaction binding the contract method 0x17ac1e6a.
//
// Solidity: function initilizeGuardians(uint256[] guardians, uint256 vote_threshold, uint256 root, address updateGuardianVerifierAddress, address socialRecoveryVerifierAddress, address poseidonContractAddress) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactor) InitilizeGuardians(opts *bind.TransactOpts, guardians []*big.Int, vote_threshold *big.Int, root *big.Int, updateGuardianVerifierAddress common.Address, socialRecoveryVerifierAddress common.Address, poseidonContractAddress common.Address) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.contract.Transact(opts, "initilizeGuardians", guardians, vote_threshold, root, updateGuardianVerifierAddress, socialRecoveryVerifierAddress, poseidonContractAddress)
}

// InitilizeGuardians is a paid mutator transaction binding the contract method 0x17ac1e6a.
//
// Solidity: function initilizeGuardians(uint256[] guardians, uint256 vote_threshold, uint256 root, address updateGuardianVerifierAddress, address socialRecoveryVerifierAddress, address poseidonContractAddress) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) InitilizeGuardians(guardians []*big.Int, vote_threshold *big.Int, root *big.Int, updateGuardianVerifierAddress common.Address, socialRecoveryVerifierAddress common.Address, poseidonContractAddress common.Address) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.InitilizeGuardians(&_PrivateRecoveryAccount.TransactOpts, guardians, vote_threshold, root, updateGuardianVerifierAddress, socialRecoveryVerifierAddress, poseidonContractAddress)
}

// InitilizeGuardians is a paid mutator transaction binding the contract method 0x17ac1e6a.
//
// Solidity: function initilizeGuardians(uint256[] guardians, uint256 vote_threshold, uint256 root, address updateGuardianVerifierAddress, address socialRecoveryVerifierAddress, address poseidonContractAddress) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorSession) InitilizeGuardians(guardians []*big.Int, vote_threshold *big.Int, root *big.Int, updateGuardianVerifierAddress common.Address, socialRecoveryVerifierAddress common.Address, poseidonContractAddress common.Address) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.InitilizeGuardians(&_PrivateRecoveryAccount.TransactOpts, guardians, vote_threshold, root, updateGuardianVerifierAddress, socialRecoveryVerifierAddress, poseidonContractAddress)
}

// Recover is a paid mutator transaction binding the contract method 0x8af96f19.
//
// Solidity: function recover(address newOwner, uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[3] input) returns(bool valid, bool update)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactor) Recover(opts *bind.TransactOpts, newOwner common.Address, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [3]*big.Int) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.contract.Transact(opts, "recover", newOwner, a, b, c, input)
}

// Recover is a paid mutator transaction binding the contract method 0x8af96f19.
//
// Solidity: function recover(address newOwner, uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[3] input) returns(bool valid, bool update)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) Recover(newOwner common.Address, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [3]*big.Int) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.Recover(&_PrivateRecoveryAccount.TransactOpts, newOwner, a, b, c, input)
}

// Recover is a paid mutator transaction binding the contract method 0x8af96f19.
//
// Solidity: function recover(address newOwner, uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[3] input) returns(bool valid, bool update)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorSession) Recover(newOwner common.Address, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [3]*big.Int) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.Recover(&_PrivateRecoveryAccount.TransactOpts, newOwner, a, b, c, input)
}

// UpdateGuardian is a paid mutator transaction binding the contract method 0xd0111bd7.
//
// Solidity: function updateGuardian(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[6] input) returns(bool)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactor) UpdateGuardian(opts *bind.TransactOpts, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [6]*big.Int) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.contract.Transact(opts, "updateGuardian", a, b, c, input)
}

// UpdateGuardian is a paid mutator transaction binding the contract method 0xd0111bd7.
//
// Solidity: function updateGuardian(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[6] input) returns(bool)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) UpdateGuardian(a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [6]*big.Int) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.UpdateGuardian(&_PrivateRecoveryAccount.TransactOpts, a, b, c, input)
}

// UpdateGuardian is a paid mutator transaction binding the contract method 0xd0111bd7.
//
// Solidity: function updateGuardian(uint256[2] a, uint256[2][2] b, uint256[2] c, uint256[6] input) returns(bool)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorSession) UpdateGuardian(a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int, input [6]*big.Int) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.UpdateGuardian(&_PrivateRecoveryAccount.TransactOpts, a, b, c, input)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.UpgradeTo(&_PrivateRecoveryAccount.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.UpgradeTo(&_PrivateRecoveryAccount.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.UpgradeToAndCall(&_PrivateRecoveryAccount.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.UpgradeToAndCall(&_PrivateRecoveryAccount.TransactOpts, newImplementation, data)
}

// ValidateUserOp is a paid mutator transaction binding the contract method 0x3a871cdd.
//
// Solidity: function validateUserOp((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes) userOp, bytes32 userOpHash, uint256 missingAccountFunds) returns(uint256 validationData)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactor) ValidateUserOp(opts *bind.TransactOpts, userOp UserOperation, userOpHash [32]byte, missingAccountFunds *big.Int) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.contract.Transact(opts, "validateUserOp", userOp, userOpHash, missingAccountFunds)
}

// ValidateUserOp is a paid mutator transaction binding the contract method 0x3a871cdd.
//
// Solidity: function validateUserOp((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes) userOp, bytes32 userOpHash, uint256 missingAccountFunds) returns(uint256 validationData)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) ValidateUserOp(userOp UserOperation, userOpHash [32]byte, missingAccountFunds *big.Int) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.ValidateUserOp(&_PrivateRecoveryAccount.TransactOpts, userOp, userOpHash, missingAccountFunds)
}

// ValidateUserOp is a paid mutator transaction binding the contract method 0x3a871cdd.
//
// Solidity: function validateUserOp((address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes) userOp, bytes32 userOpHash, uint256 missingAccountFunds) returns(uint256 validationData)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorSession) ValidateUserOp(userOp UserOperation, userOpHash [32]byte, missingAccountFunds *big.Int) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.ValidateUserOp(&_PrivateRecoveryAccount.TransactOpts, userOp, userOpHash, missingAccountFunds)
}

// WithdrawDepositTo is a paid mutator transaction binding the contract method 0x4d44560d.
//
// Solidity: function withdrawDepositTo(address withdrawAddress, uint256 amount) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactor) WithdrawDepositTo(opts *bind.TransactOpts, withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.contract.Transact(opts, "withdrawDepositTo", withdrawAddress, amount)
}

// WithdrawDepositTo is a paid mutator transaction binding the contract method 0x4d44560d.
//
// Solidity: function withdrawDepositTo(address withdrawAddress, uint256 amount) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) WithdrawDepositTo(withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.WithdrawDepositTo(&_PrivateRecoveryAccount.TransactOpts, withdrawAddress, amount)
}

// WithdrawDepositTo is a paid mutator transaction binding the contract method 0x4d44560d.
//
// Solidity: function withdrawDepositTo(address withdrawAddress, uint256 amount) returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorSession) WithdrawDepositTo(withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.WithdrawDepositTo(&_PrivateRecoveryAccount.TransactOpts, withdrawAddress, amount)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PrivateRecoveryAccount.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountSession) Receive() (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.Receive(&_PrivateRecoveryAccount.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_PrivateRecoveryAccount *PrivateRecoveryAccountTransactorSession) Receive() (*types.Transaction, error) {
	return _PrivateRecoveryAccount.Contract.Receive(&_PrivateRecoveryAccount.TransactOpts)
}

// PrivateRecoveryAccountAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the PrivateRecoveryAccount contract.
type PrivateRecoveryAccountAdminChangedIterator struct {
	Event *PrivateRecoveryAccountAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PrivateRecoveryAccountAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrivateRecoveryAccountAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PrivateRecoveryAccountAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PrivateRecoveryAccountAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrivateRecoveryAccountAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrivateRecoveryAccountAdminChanged represents a AdminChanged event raised by the PrivateRecoveryAccount contract.
type PrivateRecoveryAccountAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*PrivateRecoveryAccountAdminChangedIterator, error) {

	logs, sub, err := _PrivateRecoveryAccount.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &PrivateRecoveryAccountAdminChangedIterator{contract: _PrivateRecoveryAccount.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *PrivateRecoveryAccountAdminChanged) (event.Subscription, error) {

	logs, sub, err := _PrivateRecoveryAccount.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrivateRecoveryAccountAdminChanged)
				if err := _PrivateRecoveryAccount.contract.UnpackLog(event, "AdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) ParseAdminChanged(log types.Log) (*PrivateRecoveryAccountAdminChanged, error) {
	event := new(PrivateRecoveryAccountAdminChanged)
	if err := _PrivateRecoveryAccount.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrivateRecoveryAccountBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the PrivateRecoveryAccount contract.
type PrivateRecoveryAccountBeaconUpgradedIterator struct {
	Event *PrivateRecoveryAccountBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PrivateRecoveryAccountBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrivateRecoveryAccountBeaconUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PrivateRecoveryAccountBeaconUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PrivateRecoveryAccountBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrivateRecoveryAccountBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrivateRecoveryAccountBeaconUpgraded represents a BeaconUpgraded event raised by the PrivateRecoveryAccount contract.
type PrivateRecoveryAccountBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*PrivateRecoveryAccountBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _PrivateRecoveryAccount.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &PrivateRecoveryAccountBeaconUpgradedIterator{contract: _PrivateRecoveryAccount.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *PrivateRecoveryAccountBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _PrivateRecoveryAccount.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrivateRecoveryAccountBeaconUpgraded)
				if err := _PrivateRecoveryAccount.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) ParseBeaconUpgraded(log types.Log) (*PrivateRecoveryAccountBeaconUpgraded, error) {
	event := new(PrivateRecoveryAccountBeaconUpgraded)
	if err := _PrivateRecoveryAccount.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrivateRecoveryAccountInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the PrivateRecoveryAccount contract.
type PrivateRecoveryAccountInitializedIterator struct {
	Event *PrivateRecoveryAccountInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PrivateRecoveryAccountInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrivateRecoveryAccountInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PrivateRecoveryAccountInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PrivateRecoveryAccountInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrivateRecoveryAccountInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrivateRecoveryAccountInitialized represents a Initialized event raised by the PrivateRecoveryAccount contract.
type PrivateRecoveryAccountInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) FilterInitialized(opts *bind.FilterOpts) (*PrivateRecoveryAccountInitializedIterator, error) {

	logs, sub, err := _PrivateRecoveryAccount.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PrivateRecoveryAccountInitializedIterator{contract: _PrivateRecoveryAccount.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PrivateRecoveryAccountInitialized) (event.Subscription, error) {

	logs, sub, err := _PrivateRecoveryAccount.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrivateRecoveryAccountInitialized)
				if err := _PrivateRecoveryAccount.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) ParseInitialized(log types.Log) (*PrivateRecoveryAccountInitialized, error) {
	event := new(PrivateRecoveryAccountInitialized)
	if err := _PrivateRecoveryAccount.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrivateRecoveryAccountPrivateRecoveryAccountInitializedIterator is returned from FilterPrivateRecoveryAccountInitialized and is used to iterate over the raw logs and unpacked data for PrivateRecoveryAccountInitialized events raised by the PrivateRecoveryAccount contract.
type PrivateRecoveryAccountPrivateRecoveryAccountInitializedIterator struct {
	Event *PrivateRecoveryAccountPrivateRecoveryAccountInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PrivateRecoveryAccountPrivateRecoveryAccountInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrivateRecoveryAccountPrivateRecoveryAccountInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PrivateRecoveryAccountPrivateRecoveryAccountInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PrivateRecoveryAccountPrivateRecoveryAccountInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrivateRecoveryAccountPrivateRecoveryAccountInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrivateRecoveryAccountPrivateRecoveryAccountInitialized represents a PrivateRecoveryAccountInitialized event raised by the PrivateRecoveryAccount contract.
type PrivateRecoveryAccountPrivateRecoveryAccountInitialized struct {
	EntryPoint common.Address
	Owner      common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPrivateRecoveryAccountInitialized is a free log retrieval operation binding the contract event 0xb89fd6570b8a0665e4930593cdd2e42ae81785b67e283f02bac6079b176ad409.
//
// Solidity: event PrivateRecoveryAccountInitialized(address indexed entryPoint, address indexed owner)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) FilterPrivateRecoveryAccountInitialized(opts *bind.FilterOpts, entryPoint []common.Address, owner []common.Address) (*PrivateRecoveryAccountPrivateRecoveryAccountInitializedIterator, error) {

	var entryPointRule []interface{}
	for _, entryPointItem := range entryPoint {
		entryPointRule = append(entryPointRule, entryPointItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _PrivateRecoveryAccount.contract.FilterLogs(opts, "PrivateRecoveryAccountInitialized", entryPointRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &PrivateRecoveryAccountPrivateRecoveryAccountInitializedIterator{contract: _PrivateRecoveryAccount.contract, event: "PrivateRecoveryAccountInitialized", logs: logs, sub: sub}, nil
}

// WatchPrivateRecoveryAccountInitialized is a free log subscription operation binding the contract event 0xb89fd6570b8a0665e4930593cdd2e42ae81785b67e283f02bac6079b176ad409.
//
// Solidity: event PrivateRecoveryAccountInitialized(address indexed entryPoint, address indexed owner)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) WatchPrivateRecoveryAccountInitialized(opts *bind.WatchOpts, sink chan<- *PrivateRecoveryAccountPrivateRecoveryAccountInitialized, entryPoint []common.Address, owner []common.Address) (event.Subscription, error) {

	var entryPointRule []interface{}
	for _, entryPointItem := range entryPoint {
		entryPointRule = append(entryPointRule, entryPointItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _PrivateRecoveryAccount.contract.WatchLogs(opts, "PrivateRecoveryAccountInitialized", entryPointRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrivateRecoveryAccountPrivateRecoveryAccountInitialized)
				if err := _PrivateRecoveryAccount.contract.UnpackLog(event, "PrivateRecoveryAccountInitialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePrivateRecoveryAccountInitialized is a log parse operation binding the contract event 0xb89fd6570b8a0665e4930593cdd2e42ae81785b67e283f02bac6079b176ad409.
//
// Solidity: event PrivateRecoveryAccountInitialized(address indexed entryPoint, address indexed owner)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) ParsePrivateRecoveryAccountInitialized(log types.Log) (*PrivateRecoveryAccountPrivateRecoveryAccountInitialized, error) {
	event := new(PrivateRecoveryAccountPrivateRecoveryAccountInitialized)
	if err := _PrivateRecoveryAccount.contract.UnpackLog(event, "PrivateRecoveryAccountInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrivateRecoveryAccountUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the PrivateRecoveryAccount contract.
type PrivateRecoveryAccountUpgradedIterator struct {
	Event *PrivateRecoveryAccountUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PrivateRecoveryAccountUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrivateRecoveryAccountUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PrivateRecoveryAccountUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PrivateRecoveryAccountUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrivateRecoveryAccountUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrivateRecoveryAccountUpgraded represents a Upgraded event raised by the PrivateRecoveryAccount contract.
type PrivateRecoveryAccountUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*PrivateRecoveryAccountUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PrivateRecoveryAccount.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &PrivateRecoveryAccountUpgradedIterator{contract: _PrivateRecoveryAccount.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *PrivateRecoveryAccountUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _PrivateRecoveryAccount.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrivateRecoveryAccountUpgraded)
				if err := _PrivateRecoveryAccount.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_PrivateRecoveryAccount *PrivateRecoveryAccountFilterer) ParseUpgraded(log types.Log) (*PrivateRecoveryAccountUpgraded, error) {
	event := new(PrivateRecoveryAccountUpgraded)
	if err := _PrivateRecoveryAccount.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
