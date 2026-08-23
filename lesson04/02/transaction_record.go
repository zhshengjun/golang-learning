package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

func Record(txHash string) {
	rawUrl := viper.GetString("raw.url") // RPC 地址

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	client, _ := ethclient.DialContext(ctx, rawUrl)
	defer client.Close()

	tx, _, _ := client.TransactionByHash(context.Background(), common.HexToHash(txHash))

	from, _ := types.Sender(
		types.LatestSignerForChainID(tx.ChainId()),
		tx,
	)
	fmt.Println("=================Transaction Part=========================")
	fmt.Printf("Transaction From: %s\n", from.String())
	fmt.Printf("Transaction To: %s\n", tx.To().String())
	fmt.Printf("Transaction Nonce: %d\n", tx.Nonce())
	fmt.Printf("Transaction Time: %s\n", tx.Time().Format(time.DateTime))

	fmt.Printf("Transaction Amount（wei）: %s\n", tx.Value())  // 这是交易金额
	fmt.Printf("Transaction Input: %s\n", string(tx.Data())) // 这是输入内容

	fmt.Printf("Gas Fee Cap: %s\n", tx.GasFeeCap().String())
	fmt.Printf("Gas Price: %s\n", tx.GasPrice().String())
	fmt.Printf("Gas Tip Cap: %s\n", tx.GasTipCap().String())
	fmt.Printf("Gas TipCap Cmp: %d\n", tx.GasTipCapCmp)
	fmt.Printf("Cost = Amount + GasFeeCap * GasLimit: %s\n", tx.Cost().String()) //实时计算的理论

	fmt.Println("=================Receipt Part=========================")

	receipt, _ := client.TransactionReceipt(ctx, common.HexToHash(txHash))
	fmt.Printf("Receipt Status: %d\n", receipt.Status) // 这是交易结果
	fmt.Printf("Receipt BlockNumber: %s\n", receipt.BlockNumber)
	fmt.Printf("Receipt BlockHash: %s\n", receipt.BlockHash)

	fmt.Printf("Receipt ContractAddress: %s\n", receipt.ContractAddress.Hex()) // 这里没有合约 现在是 0x0000000000000000000000000000000000000000

	fmt.Printf("Receipt GasUsed: %d\n", receipt.GasUsed)
	fmt.Printf("Receipt EffectiveGasPrice: %d\n", receipt.EffectiveGasPrice)

	gasPrice := new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), receipt.EffectiveGasPrice)
	fmt.Printf("GasUsed * EffectiveGasPrice(wei): %d\n", gasPrice)                                                         // 这是本次交易，设计的gas费用
	fmt.Printf("GasUsed * EffectiveGasPrice(ETH): %s\n", new(big.Rat).SetFrac(gasPrice, big.NewInt(1e18)).FloatString(18)) // 这是本次交易，设计的gas费用

	// 普通的转账是没有的
	if receipt.BlobGasUsed > 0 {
		fmt.Printf("Receipt BlobGasUsed: %d\n", receipt.BlobGasUsed)
		fmt.Printf("Receipt BlobGasPrice: %d\n", receipt.BlobGasPrice)
		fmt.Printf("Blob GasUsed : %d\n", new(big.Int).Mul(new(big.Int).SetUint64(receipt.BlobGasUsed), receipt.BlobGasPrice))
	}
}
