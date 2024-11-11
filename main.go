package main

import (
	"log"
	"os"

	"blockchain-parser/parser"
)

func main() {
	// Check if an Ethereum address was provided as an argument
	if len(os.Args) < 2 {
		log.Fatal("Please provide an Ethereum address as an argument")
	}

	ethAddress := os.Args[1] // Get the Ethereum address from the first argument

	ethParser := parser.NewEthereumParser()
	if !ethParser.Subscribe(ethAddress) {
		log.Fatalf("Failed to subscribe to the Ethereum address: %s", ethAddress)
	}

	log.Printf("Subscribed to Ethereum address: %s", ethAddress)
	log.Println("Starting the Ethereum block parser...")
	parser.StartParser(ethParser)
}
