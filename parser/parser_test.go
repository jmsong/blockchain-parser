package parser

import (
	"blockchain-parser/models"
	"testing"
)

func TestNewEthereumParser(t *testing.T) {
	parser := NewEthereumParser()
	if parser == nil {
		t.Fatal("Expected non-nil parser instance")
	}
	if len(parser.subscribed) != 0 {
		t.Errorf("Expected empty subscribed map, got %d elements", len(parser.subscribed))
	}
}

func TestSubscribe(t *testing.T) {
	parser := NewEthereumParser()

	address := "0xTestAddress"
	success := parser.Subscribe(address)
	if !success {
		t.Error("Expected subscription to return true on first call")
	}

	// Test duplicate subscription
	success = parser.Subscribe(address)
	if success {
		t.Error("Expected subscription to return false on duplicate subscription")
	}
}

func TestGetTransactions(t *testing.T) {
	parser := NewEthereumParser()
	address := "0xTestAddress"

	// Prepopulate some test transactions
	parser.mu.Lock()
	parser.transactions[address] = []models.Transaction{
		{Hash: "0xHash1", From: "0xFrom1", To: address, Value: "100"},
		{Hash: "0xHash2", From: address, To: "0xTo2", Value: "200"},
	}
	parser.mu.Unlock()

	transactions := parser.GetTransactions(address)
	if len(transactions) != 2 {
		t.Fatalf("Expected 2 transactions, got %d", len(transactions))
	}

	if transactions[0].Hash != "0xHash1" || transactions[1].Hash != "0xHash2" {
		t.Error("Transaction data does not match expected values")
	}
}
