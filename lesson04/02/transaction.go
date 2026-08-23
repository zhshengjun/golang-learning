package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

func SendTransaction(private string, to string, amount float64, input []byte) {
	rawUrl := viper.GetString("raw.url") // RPC 地址
	toAddress := common.HexToAddress(to) // 接收方

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second) // 构建上下文
	defer cancelFunc()

	ethClient, err := ethclient.DialContext(ctx, rawUrl) // 构建 ethereum 的客户端
	if err != nil {
		log.Fatal(err)
	}
	defer ethClient.Close() // 结束前关闭

	privateKey, err := crypto.HexToECDSA(private)
	publicKey := privateKey.PublicKey
	from := crypto.PubkeyToAddress(publicKey)

	fmt.Println("From Address Key:", from)
	chainID, err := ethClient.ChainID(ctx)
	fmt.Println("Chain ID:", chainID)

	nonce, err := ethClient.PendingNonceAt(ctx, from) // 获取账户的 nonce

	gasTipCap, err := ethClient.SuggestGasTipCap(ctx) // 建议的gas费用

	header, err := ethClient.HeaderByNumber(ctx, nil) // 获取Header

	baseFee := header.BaseFee
	if baseFee == nil {
		suggestGasPrice, _ := ethClient.SuggestGasPrice(ctx)

		baseFee = suggestGasPrice
	}

	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(baseFee, big.NewInt(2)),
		gasTipCap,
	)

	amountWei := new(big.Float).Mul(
		big.NewFloat(amount),
		big.NewFloat(1e18),
	)
	valueWei, _ := amountWei.Int(new(big.Int))
	gasLimit := uint64(21000)
	if input != nil {
		gasLimit, _ = ethClient.EstimateGas(ctx, ethereum.CallMsg{
			From:  from,
			To:    &toAddress,
			Value: valueWei,
			Data:  input,
		})
	}
	// 当前交易的总费用
	totalCost := new(big.Int).Add(valueWei, new(big.Int).Mul(gasFeeCap, big.NewInt(int64(gasLimit))))

	balanceAt, err := ethClient.BalanceAt(ctx, from, nil) // 余额
	if err != nil {
		log.Fatal(err)
	}

	if balanceAt.Cmp(totalCost) < 0 {
		fmt.Println("account balance is less than total cost")
		return
	}

	// 构建交易体
	txData := types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     valueWei,
		Data:      []byte("test transaction"),
	}

	// 获取交易
	tx := types.NewTx(&txData)

	// 获取签名
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		log.Fatal("SignTx:", err)
	}

	// 发送交易
	err = ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal("Send tx:", err)
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Printf("From: %s\n", from.Hex())                                                    // 0xb58bbA2158cD9E2d52985D21e863217941600734
	fmt.Printf("To: %s\n", toAddress.Hex())                                                 //0x44b29771acb6144dDbDa47A075ed64D68Fa53c1D
	fmt.Printf("Value : %s ETH(%s wei))\n", fmt.Sprintf("%.6f", amount), valueWei.String()) // 0.400000 ETH(400000000000000000 wei))
	fmt.Printf("Gas Limit: %d\n", gasLimit)                                                 // 21806
	fmt.Printf("Gas Tip Cap: %d\n", gasTipCap)                                              // 1000000
	fmt.Printf("Gas Fee Cap: %d\n", gasFeeCap)                                              // 2165191302
	fmt.Printf("Total Cost: %d\n", totalCost)                                               // 400047214161531412
	fmt.Printf("Nonce: %d\n", nonce)                                                        // 20
	fmt.Printf("Tx Hash: %s\n", signedTx.Hash().Hex())                                      // 0x81093adda2c955ee69e2bf995660006fbcfb1b71b60e2de86bd3741417a3d22b
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++")
}
