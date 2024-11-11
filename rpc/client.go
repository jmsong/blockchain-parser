package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var url = "https://cloudflare-eth.com" // Declare the URL as a package-level variable

// JSONRPCRequest sends a request to the Ethereum JSON-RPC API
func JSONRPCRequest(method string, params []interface{}) ([]byte, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// GetBlockNumber retrieves the current block number from the Ethereum blockchain
func GetBlockNumber() (int, error) {
	response, err := JSONRPCRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		return 0, err
	}
	var result map[string]interface{}
	json.Unmarshal(response, &result)
	blockHex := result["result"].(string)
	var blockNumber int
	fmt.Sscanf(blockHex, "0x%x", &blockNumber) // Convert hex to int
	return blockNumber, nil
}

// GetBlockByNumber retrieves the block information by block number
func GetBlockByNumber(blockNumber int) (map[string]interface{}, error) {
	blockParam := fmt.Sprintf("0x%x", blockNumber) // Convert block number to hex
	response, err := JSONRPCRequest("eth_getBlockByNumber", []interface{}{blockParam, true})
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	json.Unmarshal(response, &result)
	return result, nil
}
