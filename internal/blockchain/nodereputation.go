// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// NodeReputationABI is the input ABI used to generate the binding from.
const NodeReputationABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"gpuModel\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"vram\",\"type\":\"uint64\"}],\"name\":\"registerNode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sendHeartbeat\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// NodeReputation is an auto generated Go binding around an Ethereum contract.
type NodeReputation struct {
	NodeReputationCaller     // Read-only binding to the contract
	NodeReputationTransactor // Write-only binding to the contract
	NodeReputationFilterer   // Log filterer for contract events
}

// NodeReputationCaller is an auto generated read-only Go binding around an Ethereum contract.
type NodeReputationCaller struct {
	contract *bind.BoundContract // Generic wrapper around the contract
}

// NodeReputationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NodeReputationTransactor struct {
	contract *bind.BoundContract // Generic wrapper around the contract
}

// NodeReputationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NodeReputationFilterer struct {
	contract *bind.BoundContract // Generic wrapper around the contract
}

// NodeReputationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NodeReputationSession struct {
	Contract     *NodeReputation   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodeReputationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NodeReputationCallerSession struct {
	Contract *NodeReputationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// NodeReputationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NodeReputationTransactorSession struct {
	Contract     *NodeReputationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// NodeReputationRaw is an auto generated low-level Go binding around an Ethereum contract.
type NodeReputationRaw struct {
	Contract *NodeReputation // Generic contract binding to access the raw methods on
}

// NodeReputationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NodeReputationCallerRaw struct {
	Contract *NodeReputationCaller // Generic read-only contract binding to access the raw methods on
}

// NodeReputationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NodeReputationTransactorRaw struct {
	Contract *NodeReputationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodeReputation creates a new instance of NodeReputation, bound to a specific deployed contract.
func NewNodeReputation(address common.Address, backend bind.ContractBackend) (*NodeReputation, error) {
	contract, err := bindNodeReputation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodeReputation{NodeReputationCaller: NodeReputationCaller{contract: contract}, NodeReputationTransactor: NodeReputationTransactor{contract: contract}, NodeReputationFilterer: NodeReputationFilterer{contract: contract}}, nil
}

// NewNodeReputationCaller creates a new read-only instance of NodeReputation, bound to a specific deployed contract.
func NewNodeReputationCaller(address common.Address, caller bind.ContractCaller) (*NodeReputationCaller, error) {
	contract, err := bindNodeReputation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodeReputationCaller{contract: contract}, nil
}

// NewNodeReputationTransactor creates a new write-only instance of NodeReputation, bound to a specific deployed contract.
func NewNodeReputationTransactor(address common.Address, transactor bind.ContractTransactor) (*NodeReputationTransactor, error) {
	contract, err := bindNodeReputation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodeReputationTransactor{contract: contract}, nil
}

// NewNodeReputationFilterer creates a new log filterer instance of NodeReputation, bound to a specific deployed contract.
func NewNodeReputationFilterer(address common.Address, filterer bind.ContractFilterer) (*NodeReputationFilterer, error) {
	contract, err := bindNodeReputation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodeReputationFilterer{contract: contract}, nil
}

// bindNodeReputation binds a generic wrapper to an already deployed contract.
func bindNodeReputation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NodeReputationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeReputation *NodeReputationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeReputation.Contract.NodeReputationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeReputation *NodeReputationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeReputation.Contract.NodeReputationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeReputation *NodeReputationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeReputation.Contract.NodeReputationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeReputation *NodeReputationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeReputation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeReputation *NodeReputationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeReputation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeReputation *NodeReputationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeReputation.Contract.contract.Transact(opts, method, params...)
}

// RegisterNode is a paid mutator transaction binding the contract method 0x12345678.
//
// Solidity: function registerNode(string gpuModel, uint64 vram) returns()
func (_NodeReputation *NodeReputationTransactor) RegisterNode(opts *bind.TransactOpts, gpuModel string, vram *big.Int) (*types.Transaction, error) {
	return _NodeReputation.contract.Transact(opts, "registerNode", gpuModel, vram)
}

// RegisterNode is a paid mutator transaction binding the contract method 0x12345678.
//
// Solidity: function registerNode(string gpuModel, uint64 vram) returns()
func (_NodeReputation *NodeReputationSession) RegisterNode(gpuModel string, vram *big.Int) (*types.Transaction, error) {
	return _NodeReputation.Contract.RegisterNode(&_NodeReputation.TransactOpts, gpuModel, vram)
}

// RegisterNode is a paid mutator transaction binding the contract method 0x12345678.
//
// Solidity: function registerNode(string gpuModel, uint64 vram) returns()
func (_NodeReputation *NodeReputationTransactorSession) RegisterNode(gpuModel string, vram *big.Int) (*types.Transaction, error) {
	return _NodeReputation.Contract.RegisterNode(&_NodeReputation.TransactOpts, gpuModel, vram)
}

// SendHeartbeat is a paid mutator transaction binding the contract method 0x87654321.
//
// Solidity: function sendHeartbeat() returns()
func (_NodeReputation *NodeReputationTransactor) SendHeartbeat(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeReputation.contract.Transact(opts, "sendHeartbeat")
}

// SendHeartbeat is a paid mutator transaction binding the contract method 0x87654321.
//
// Solidity: function sendHeartbeat() returns()
func (_NodeReputation *NodeReputationSession) SendHeartbeat() (*types.Transaction, error) {
	return _NodeReputation.Contract.SendHeartbeat(&_NodeReputation.TransactOpts)
}

// SendHeartbeat is a paid mutator transaction binding the contract method 0x87654321.
//
// Solidity: function sendHeartbeat() returns()
func (_NodeReputation *NodeReputationTransactorSession) SendHeartbeat() (*types.Transaction, error) {
	return _NodeReputation.Contract.SendHeartbeat(&_NodeReputation.TransactOpts)
} 