package rpc

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSONRPCRequest(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Error reading request body: %v", err)
		}
		defer r.Body.Close()

		// Check the request body
		if !strings.Contains(string(body), `"method":"eth_blockNumber"`) {
			t.Errorf("Expected request to contain method 'eth_blockNumber', got: %s", body)
		}

		// Send a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x4b7"}`)) // Mock block number response
	}))
	defer server.Close()

	// Replace the URL with the mock server's URL
	originalURL := url
	url = server.URL

	// Run the test
	response, err := JSONRPCRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		t.Fatalf("Error in JSONRPCRequest: %v", err)
	}

	// Check the response content
	if !strings.Contains(string(response), `"result":"0x4b7"`) {
		t.Errorf("Expected response to contain result '0x4b7', got: %s", response)
	}

	// Reset the URL to the original value after the test
	url = originalURL
}

func TestGetBlockNumber(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x4b7"}`)) // Mock block number response
	}))
	defer server.Close()

	// Replace the URL with the mock server's URL
	originalURL := url
	url = server.URL

	// Run the test
	blockNumber, err := GetBlockNumber()
	if err != nil {
		t.Fatalf("Error in GetBlockNumber: %v", err)
	}

	// Check the block number
	expectedBlockNumber := 1207
	if blockNumber != expectedBlockNumber {
		t.Errorf("Expected block number %d, got %d", expectedBlockNumber, blockNumber)
	}

	// Reset the URL to the original value after the test
	url = originalURL
}

func TestGetBlockByNumber(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":{"transactions":[{"hash":"0x123"}]}}`)) // Mock block response
	}))
	defer server.Close()

	// Replace the URL with the mock server's URL
	originalURL := url
	url = server.URL

	// Run the test
	blockData, err := GetBlockByNumber(1207)
	if err != nil {
		t.Fatalf("Error in GetBlockByNumber: %v", err)
	}

	// Check the block data
	if len(blockData["result"].(map[string]interface{})["transactions"].([]interface{})) != 1 {
		t.Errorf("Expected 1 transaction, got: %v", blockData["result"])
	}

	// Reset the URL to the original value after the test
	url = originalURL
}
