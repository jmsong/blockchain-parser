package parser

import (
	"blockchain-parser/models"
	"blockchain-parser/rpc"
	"log"
	"sync"
)

// Parser defines the interface for parsing Ethereum blocks
type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []models.Transaction
}

// EthereumParser implements the Parser interface
type EthereumParser struct {
	currentBlock int
	subscribed   map[string]bool
	transactions map[string][]models.Transaction
	mu           sync.Mutex
}

// NewEthereumParser creates a new instance of the Ethereum parser
func NewEthereumParser() *EthereumParser {
	return &EthereumParser{
		currentBlock: 0,
		subscribed:   make(map[string]bool),
		transactions: make(map[string][]models.Transaction),
	}
}

// GetCurrentBlock returns the last parsed block
func (p *EthereumParser) GetCurrentBlock() int {
	return p.currentBlock
}

// Subscribe adds an address to the observer
func (p *EthereumParser) Subscribe(address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.subscribed[address]; exists {
		return false
	}
	p.subscribed[address] = true
	return true
}

// GetTransactions returns the list of transactions for a specific address
func (p *EthereumParser) GetTransactions(address string) []models.Transaction {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.transactions[address]
}

// ParseBlock processes the block transactions and stores them for subscribed addresses
func (p *EthereumParser) ParseBlock(blockNumber int) {
	blockData, err := rpc.GetBlockByNumber(blockNumber)
	if err != nil {
		log.Println("Error fetching block data:", err)
		return
	}

	transactions := blockData["result"].(map[string]interface{})["transactions"].([]interface{})
	for _, tx := range transactions {
		txMap := tx.(map[string]interface{})
		from := txMap["from"].(string)
		to := txMap["to"].(string)
		txHash := txMap["hash"].(string)
		value := txMap["value"].(string)

		p.mu.Lock()
		if p.subscribed[from] || p.subscribed[to] {
			if p.subscribed[from] {
				p.transactions[from] = append(p.transactions[from], models.Transaction{Hash: txHash, From: from, To: to, Value: value})
			}
			if p.subscribed[to] {
				p.transactions[to] = append(p.transactions[to], models.Transaction{Hash: txHash, From: from, To: to, Value: value})
			}
		}
		p.mu.Unlock()
	}
}

// StartParser runs the Ethereum block parser
func StartParser(p *EthereumParser) {
	for {
		blockNumber, err := rpc.GetBlockNumber()
		if err != nil {
			log.Println("Error fetching block number:", err)
			continue
		}
		if blockNumber > p.GetCurrentBlock() {
			p.ParseBlock(blockNumber)
			p.currentBlock = blockNumber
		}
	}
}
