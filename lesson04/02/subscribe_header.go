package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

func Subscribe() {

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	client, err := ethclient.DialContext(ctx, viper.GetString("raw.wss"))
	if err != nil {
		log.Fatal("client dial error:", err)
	}
	defer client.Close()

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Fatal("client subscribe error:", err)
	}

	for {

		select {
		case err := <-sub.Err():
			log.Printf("client subscribe error: %s\n", err)
			return
		case header := <-headers:
			if header == nil {
				continue
			}
			fmt.Printf("receiver Header %s\n", header.Number.String())
			printfHeader(header)
		case <-ctx.Done():
			return
		}

	}

}

func printfHeader(header *types.Header) {

	fmt.Printf("Header Hash: %s\n", header.Hash().Hex())
	fmt.Printf("Header Number: %d\n", header.Number)
	fmt.Printf("Header ParentHash: %s\n", header.ParentHash.Hex())
	fmt.Printf("Header GasUsed: %d\n", header.GasUsed)
	fmt.Printf("Header Time: %s\n", time.Unix(int64(header.Time), 0).Format(time.DateTime))

}
