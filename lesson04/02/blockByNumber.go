package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

const Endpoint = "https://eth-sepolia.g.alchemy.com/v2/alch_8a2eNrhkaMX56Vfy7KA3Y"

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := ethclient.DialContext(ctx, Endpoint)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum chain: %v", err)
	}
	fmt.Println("Chain ID:", chainID)
	headerNumber := big.NewInt(11527733)

	block, err := client.BlockByNumber(ctx, headerNumber)
	if err != nil {
		log.Fatalf("failed get header: %v", err)
	}
	blockHash := block.Hash() // 这个  header.Hash()
	withdrawals := block.Withdrawals()
	blockValue, _ := blockHash.Value()
	block.Header()

	fmt.Printf("RPC URL: %s\n", Endpoint)
	fmt.Printf("Chain ID:  %d\n", chainID)
	fmt.Printf("Block Number: %d\n", block.Number)
	fmt.Printf("Block Height: %d\n", block.Header().Number)
	fmt.Printf("Block Hash: %s\n", blockHash.Hex())
	fmt.Printf("Parent Hash: %s\n", block.Header().ParentHash.Hex())
	fmt.Printf("Block withdrawals: %s\n", withdrawals)
	fmt.Printf("Block Value: %s\n", blockValue)
	fmt.Printf("State Root: %s\n", block.Header().Root.Hex())
	fmt.Printf("Block Time: %s\n", time.Unix(int64(block.Header().Time), 0).Format(time.RFC3339))
}
