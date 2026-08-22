package main

import (
	"01/commons"
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := ethclient.DialContext(ctx, commons.Endpoint)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum chain: %v", err)
	}
	fmt.Println("Chain ID:", chainID)
	_ = big.NewInt(11527733)

	block, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("failed get header: %v", err)
	}

	fmt.Printf("RPC URL: %s\n", commons.Endpoint)
	fmt.Printf("Chain ID:  %d\n", chainID)
	fmt.Printf("Block Number: %s\n", block.Number().String())
	fmt.Printf("Block Hash: %s\n", block.Hash().Hex())
	fmt.Printf("Block Head Hash: %s\n", block.Header().Hash())
	fmt.Printf("Parent Hash: %s\n", block.Header().ParentHash.Hex())
	fmt.Printf("Block withdrawals: %d\n", block.Withdrawals())
	fmt.Printf("Block GasLimit: %d\n", block.GasLimit())
	fmt.Printf("Block GasUed: %d\n", block.GasUsed())
	fmt.Printf("State Root: %s\n", block.Header().Root.Hex())
	fmt.Printf("State SlotNumber: %d\n", block.Header().SlotNumber)
	fmt.Printf("Block Time: %s\n", time.Unix(int64(block.Header().Time), 0).Format(time.RFC3339))
}
