package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockchainClient defines the interface for blockchain operations
type BlockchainClient interface {
	RegisterNode(ctx context.Context, gpuModel string, vram uint64) error
	SendHeartbeat(ctx context.Context) error
}

// ethClient implements BlockchainClient using Ethereum
type ethClient struct {
	client     *ethclient.Client
	contract   *NodeReputation
	privateKey *ecdsa.PrivateKey
	address    common.Address
	chainID    *big.Int
}

// NewEthClient creates a new Ethereum blockchain client
func NewEthClient(rpcURL, privateKeyHex, contractAddress string) (BlockchainClient, error) {
	// Connect to the Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get the public key and address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get chain ID (opBNB Testnet = 5611)
	chainID := big.NewInt(5611)

	// Parse contract address
	contractAddr := common.HexToAddress(contractAddress)

	// Create contract instance
	contract, err := NewNodeReputation(contractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract instance: %w", err)
	}

	return &ethClient{
		client:     client,
		contract:   contract,
		privateKey: privateKey,
		address:    address,
		chainID:    chainID,
	}, nil
}

// RegisterNode registers the node with the smart contract
func (e *ethClient) RegisterNode(ctx context.Context, gpuModel string, vram uint64) error {
	// Create auth for transaction
	auth, err := bind.NewKeyedTransactorWithChainID(e.privateKey, e.chainID)
	if err != nil {
		return fmt.Errorf("failed to create transaction auth: %w", err)
	}

	// Convert vram to big.Int
	vramBig := new(big.Int).SetUint64(vram)

	// Call the contract
	tx, err := e.contract.RegisterNode(auth, gpuModel, vramBig)
	if err != nil {
		return fmt.Errorf("failed to register node: %w", err)
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(ctx, e.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	return nil
}

// SendHeartbeat sends a heartbeat to the smart contract
func (e *ethClient) SendHeartbeat(ctx context.Context) error {
	// Create auth for transaction
	auth, err := bind.NewKeyedTransactorWithChainID(e.privateKey, e.chainID)
	if err != nil {
		return fmt.Errorf("failed to create transaction auth: %w", err)
	}

	// Call the contract
	tx, err := e.contract.SendHeartbeat(auth)
	if err != nil {
		return fmt.Errorf("failed to send heartbeat: %w", err)
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(ctx, e.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction: %w", err)
	}

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	return nil
}
