package models

// Transaction represents a simplified Ethereum transaction
type Transaction struct {
	Hash  string
	From  string
	To    string
	Value string
}
