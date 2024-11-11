package main

import (
	"log"

	"blockchain-parser/parser"
)

func main() {
	ethParser := parser.NewEthereumParser()
	ethParser.Subscribe("0xYourEthereumAddress") // Add the Ethereum address to monitor

	log.Println("Starting the Ethereum block parser...")
	parser.StartParser(ethParser)
}
