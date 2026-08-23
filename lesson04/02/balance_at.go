package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

func BalanceAt(addr string) *big.Int {
	rawUrl := viper.GetString("raw.url") // RPC 地址

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	client, _ := ethclient.DialContext(ctx, rawUrl)
	defer client.Close()

	address := common.HexToAddress(addr)

	fmt.Println("=================Balance Part=========================")
	fmt.Printf("address: %s\n", address)

	balance, _ := client.BalanceAt(ctx, address, nil)
	fmt.Println("balance(wei)", balance.Int64())
	ethValue := new(big.Rat).SetFrac(balance, big.NewInt(1e18))
	fmt.Println("balance(ETH)", ethValue.FloatString(3))
	return balance

}
