package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	Endpoint = "https://eth-sepolia.g.alchemy.com/v2/alch_8a2eNrhkaMX56Vfy7KA3Y"
)

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
	hash := common.HexToHash("8be5ea9f3dfc46a87155114170e78a20e3d23d2bd8ba005ce49d713ccef8d8bf")

	header, err := client.HeaderByHash(ctx, hash)
	if err != nil {
		log.Fatalf("failed get header: %v", err)
	}

	fmt.Printf("RPC URL: %s\n", Endpoint)
	fmt.Printf("Chain ID:  %d\n", chainID)
	fmt.Printf("Block Number: %d\n", header.Number)
	fmt.Printf("Block Hash: %s\n", header.Hash().Hex())
	fmt.Printf("Parent Hash: %s\n", header.ParentHash.Hex())
	fmt.Printf("State Root: %s\n", header.Root.Hex())
	fmt.Printf("Block Time: %s\n", time.Unix(int64(header.Time), 0).Format(time.RFC3339))
}
