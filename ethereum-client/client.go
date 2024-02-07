package ethereum_client

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

type Node struct {
	client          *ethclient.Client
	url             string
	chain           string
	networkID       *big.Int
	chainID         *big.Int
	signer          types.Signer
	defaultGasLimit uint64
}

type NodeImpl interface {
	GetCodeAt(address string) ([]byte, error)
	GetBalance(address string) (*big.Int, error)
	GetLatestBlockNumber() (uint64, error)
	GetBlockByNumber(number *big.Int) (*types.Block, error)
	GetTxReceipt(tx string) (*types.Receipt, error)
	GetSigner() types.Signer
	GetNonce(from string) (uint64, error)
	GetFeePerGas() (*big.Int, error)
	GetBaseFee() (*big.Int, error)
	GetChainID() *big.Int
	GetDefaultGasLimit() uint64
	SendTransaction(tx *types.Transaction) error
}

func NewNode(dialUrl, chain string, defaultGasLimit uint64) (NodeImpl, error) {
	n := &Node{
		url:             dialUrl,
		chain:           chain,
		defaultGasLimit: defaultGasLimit,
	}
	var err error

	if n.client, err = ethclient.Dial(dialUrl); err != nil {
		return nil, err
	} else if n.chainID, err = n.client.ChainID(context.Background()); err != nil {
		return nil, err
	} else if n.networkID, err = n.client.NetworkID(context.Background()); err != nil {
		return nil, err
	} else {
		n.signer = types.NewLondonSigner(n.chainID)
		log.Println("Success To Connect Node", "URL", dialUrl, "Chain", chain)
	}

	return n, nil
}

func (n *Node) getClient() (*ethclient.Client, error) {
	if n.client == nil {
		return nil, errors.New("Client is nil")
	} else {
		return n.client, nil
	}
}

func (n *Node) getPendingNonce(from string) (uint64, error) {
	if client, err := n.getClient(); err != nil {
		return 0, err
	} else if res, err := client.PendingNonceAt(context.Background(), toAddress(from)); err != nil {
		return 0, err
	} else {
		return res, nil
	}
}

func (n *Node) SendTransaction(tx *types.Transaction) error {
	if client, err := n.getClient(); err != nil {
		return err
	} else if err := client.SendTransaction(context.Background(), tx); err != nil {
		return err
	} else {
		return nil
	}
}

func (n *Node) GetNonce(from string) (uint64, error) {
	if nonce, err := n.getPendingNonce(from); err != nil {
		return 0, err
	} else {
		return nonce, nil
	}
}

func (n *Node) GetSigner() types.Signer {
	return n.signer
}

func (n *Node) GetBalance(address string) (*big.Int, error) {
	if client, err := n.getClient(); err != nil {
		return nil, err
	} else if balance, err := client.BalanceAt(context.Background(), toAddress(address), nil); err != nil {
		return nil, err
	} else {
		return balance, nil
	}
}

func (n *Node) GetLatestBlockNumber() (uint64, error) {
	if client, err := n.getClient(); err != nil {
		return 0, err
	} else if num, err := client.BlockNumber(context.Background()); err != nil {
		return 0, err
	} else {
		return num, nil
	}
}

func (n *Node) GetBlockByNumber(number *big.Int) (*types.Block, error) {
	if client, err := n.getClient(); err != nil {
		return nil, err
	} else if res, err := client.BlockByNumber(context.Background(), number); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (n *Node) GetTxReceipt(tx string) (*types.Receipt, error) {
	if client, err := n.getClient(); err != nil {
		return nil, err
	} else if res, err := client.TransactionReceipt(context.Background(), toHash(tx)); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (n *Node) GetCodeAt(address string) ([]byte, error) {
	if client, err := n.getClient(); err != nil {
		return nil, err
	} else if bytes, err := client.CodeAt(context.Background(), toAddress(address), nil); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

func (n *Node) GetFeePerGas() (*big.Int, error) {
	if client, err := n.getClient(); err != nil {
		return nil, err
	} else if res, err := client.SuggestGasTipCap(context.Background()); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (n *Node) GetChainID() *big.Int {
	return n.chainID
}

func (n *Node) GetBaseFee() (*big.Int, error) {
	if number, err := n.GetLatestBlockNumber(); err != nil {
		return nil, err
	} else if block, err := n.GetBlockByNumber(big.NewInt(int64(number))); err != nil {
		return nil, err
	} else {
		return block.BaseFee(), nil
	}
}

func (n *Node) GetDefaultGasLimit() uint64 {
	return n.defaultGasLimit
}

func toAddress(str string) common.Address {
	return common.HexToAddress(str)
}

func toHash(str string) common.Hash {
	return common.HexToHash(str)
}
