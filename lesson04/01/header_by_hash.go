package main

import (
	"01/commons"
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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

	latestHeader, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("failed get header: %v", err)
	}

	fmt.Printf("Latest Block Number: %d\n", latestHeader.Number)
	fmt.Printf("Latest Block Hash: %s\n", latestHeader.Hash().Hex())
	fmt.Printf("Latest Parent Hash: %s\n", latestHeader.ParentHash.Hex())
	fmt.Printf("Latest State Root: %s\n", latestHeader.Root.Hex())
	fmt.Printf("Latest Block Time: %s\n", time.Unix(int64(latestHeader.Time), 0).Format(time.RFC3339))

	safeHeader, err := client.HeaderByNumber(
		ctx,
		big.NewInt(int64(rpc.SafeBlockNumber)), // -4 -> "safe"
	)

	fmt.Printf("Safe Block Number: %d\n", safeHeader.Number)
	fmt.Printf("Safe Block Hash: %s\n", safeHeader.Hash().Hex())
	fmt.Printf("Safe Parent Hash: %s\n", safeHeader.ParentHash.Hex())
	fmt.Printf("Safe State Root: %s\n", safeHeader.Root.Hex())
	fmt.Printf("Safe Block Time: %s\n", time.Unix(int64(safeHeader.Time), 0).Format(time.RFC3339))
	fmt.Printf("Safe VS Latest Confirmations Number:: %d\n", safeHeader.Number.Int64()-latestHeader.Number.Int64())

	finalizedHeader, err := client.HeaderByNumber(
		ctx,
		big.NewInt(int64(rpc.FinalizedBlockNumber)), // -3 -> "finalized"
	)
	fmt.Printf("Finalized Block Number: %d\n", finalizedHeader.Number)
	fmt.Printf("Finalized Block Hash: %s\n", finalizedHeader.Hash().Hex())
	fmt.Printf("Finalized Parent Hash: %s\n", finalizedHeader.ParentHash.Hex())
	fmt.Printf("Finalized State Root: %s\n", finalizedHeader.Root.Hex())
	fmt.Printf("Finalized Block Time: %s\n", time.Unix(int64(finalizedHeader.Time), 0).Format(time.RFC3339))
	fmt.Printf("Finalized VS Latest Confirmations Number:: %d\n", finalizedHeader.Number.Int64()-safeHeader.Number.Int64())
}
