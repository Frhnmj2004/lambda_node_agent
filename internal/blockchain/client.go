package blockchain

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockchainClient defines the interface for blockchain operations required by the agent.
type BlockchainClient interface {
	// RegisterNode registers the node on-chain with its GPU model and VRAM.
	RegisterNode(ctx context.Context, gpuModel string, vram uint64) (string, error)
	// SendHeartbeat sends a heartbeat transaction to the contract.
	SendHeartbeat(ctx context.Context) (string, error)
}

// ethClient implements BlockchainClient using go-ethereum.
type ethClient struct {
	client       *ethclient.Client
	transactor   *bind.TransactOpts
	contractAddr common.Address
	// contract   *NodeReputation // TODO: Replace with actual contract binding
}

// NewEthClient creates a new ethClient instance.
func NewEthClient(rpcUrl, privateKey, contractAddr string) (BlockchainClient, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	pk, err := crypto.HexToECDSA(trimHexPrefix(privateKey))
	if err != nil {
		return nil, err
	}

	fromAddr := crypto.PubkeyToAddress(pk.PublicKey)
	transactor, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(5611)) // TODO: Make chain ID configurable
	if err != nil {
		return nil, err
	}
	transactor.From = fromAddr

	contractAddress := common.HexToAddress(contractAddr)

	// TODO: Initialize contract binding with contractAddress and client

	return &ethClient{
		client:       client,
		transactor:   transactor,
		contractAddr: contractAddress,
		// contract:  contract, // TODO
	}, nil
}

// RegisterNode registers the node on-chain with its GPU model and VRAM.
func (e *ethClient) RegisterNode(ctx context.Context, gpuModel string, vram uint64) (string, error) {
	// TODO: Call contract's registerNode method
	// Example: tx, err := e.contract.RegisterNode(e.transactor, gpuModel, vram)
	return "", errors.New("RegisterNode not implemented: contract binding required")
}

// SendHeartbeat sends a heartbeat transaction to the contract.
func (e *ethClient) SendHeartbeat(ctx context.Context) (string, error) {
	// TODO: Call contract's sendHeartbeat method
	// Example: tx, err := e.contract.SendHeartbeat(e.transactor)
	return "", errors.New("SendHeartbeat not implemented: contract binding required")
}

// Helper to trim 0x prefix from hex string
func trimHexPrefix(s string) string {
	if len(s) >= 2 && s[:2] == "0x" {
		return s[2:]
	}
	return s
}

// NodeReputation is a placeholder for the contract binding. Replace with actual generated Go binding.
type NodeReputation struct{}
