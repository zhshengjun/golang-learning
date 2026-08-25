package transaction

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
)

// TransactionOriginal 操作原生代币
func TransactionOriginal(amount float64) {

	//  查询余额
	BalanceAtByPrivateKey(os.Getenv("ACCOUNT_FROM_PRIVATE")) //根据私钥获取

	// 获取区块信息
	BlockInfo(big.NewInt(int64(rpc.LatestBlockNumber)))

	// 模拟交易
	Transaction(os.Getenv("ACCOUNT_FROM_PRIVATE"), os.Getenv("ACCOUNT_TO_ADDRESS"), amount)
}

// Transaction 模拟交易
func Transaction(privateHex string, toAddressHex string, amount float64) {

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()
	client, err := ethclient.DialContext(ctx, viper.GetString("raw.url"))
	if err != nil {
		log.Fatalf("failed to connect to raw client: %v", err)
	}
	defer client.Close()

	privateKey, _ := crypto.HexToECDSA(privateHex)
	form := Private2Address(privateHex)
	to := common.HexToAddress(toAddressHex)
	nonce, _ := client.PendingNonceAt(ctx, form)

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get header: %v", err)
	}

	baseFee := header.BaseFee
	// 若不支持EIP-1559
	if baseFee == nil {
		suggestGasPrice, _ := client.SuggestGasPrice(ctx)
		baseFee = suggestGasPrice
	}

	gasTipCap, _ := client.SuggestGasTipCap(ctx)

	gasFeeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), gasTipCap)

	// 这里计算wei，注意amount不能直接转int，容易丢数据
	value, _ := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e18)).Int(new(big.Int))

	// 尝试计算
	gasLimit, _ := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  form,
		To:    &to,
		Value: value,
		Data:  nil,
	})

	maxCost := new(big.Int).Add(value, new(big.Int).Mul(gasFeeCap, big.NewInt(int64(gasLimit))))

	balanceAt, _ := client.BalanceAt(ctx, form, nil)

	if balanceAt.Cmp(maxCost) < 0 {
		log.Fatalf("Balance is not enough to pay for %s", to)
	}

	chainID, _ := client.ChainID(ctx)
	// 构造交易体
	txData := types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &to,
		Value:     value,
		Data:      nil,
	}

	tx := types.NewTx(&txData)
	// 签名 ,注意 signer需要匹配
	transaction, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		log.Fatalf("failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(ctx, transaction)
	if err != nil {
		log.Fatalf("failed to send transaction: %v", err)
	}
	fmt.Printf("transaction sent successfully")
	fmt.Printf("Tx Hash :%s\n", transaction.Hash().Hex())
}

// BlockInfo 获取Block 的信息
func BlockInfo(blockNumber *big.Int) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()
	client, err := ethclient.DialContext(ctx, os.Getenv("RAW_URL"))
	if err != nil {
		log.Fatalf("failed to connect to raw client: %v", err)
	}
	defer client.Close()

	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		log.Fatalf("failed to get header by block number: %v", err)
	}
	fmt.Println("=================Block Info============================")
	fmt.Printf("Block Number: %d\n", block.Number().Int64())
	fmt.Printf("Block Hash: %s\n", block.Hash().Hex())
	fmt.Printf("Block ParentHash: %s\n", block.ParentHash().Hex())
	fmt.Printf("Block TxHash: %s\n", block.TxHash().Hex())
	fmt.Printf("Block Transactions Number: %d\n", len(block.Transactions()))
	fmt.Printf("Block GasUsed: %d(wei)\n", block.GasUsed())
	fmt.Printf("Block Time: %s\n", time.Unix(int64(block.Time()), 0).Format(time.RFC3339))
}

// BalanceAtByPrivateKey 根据私钥获取余额
func BalanceAtByPrivateKey(privateKey string) {
	//1. 私钥解析
	address := Private2Address(privateKey)

	BalanceAtByAddress(address.Hex()) // 根据地址获取余额
}

func Private2Address(privateKey string) common.Address {
	privateToECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateToECDSA.PublicKey

	address := crypto.PubkeyToAddress(publicKey)

	fmt.Printf("privateKey to address: %v\n", address.Hex())
	return address
}

// BalanceAtByAddress 根据地址获取余额
func BalanceAtByAddress(addressStr string) {

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()
	client, err := ethclient.DialContext(ctx, os.Getenv("RAW_URL"))
	if err != nil {
		log.Fatalf("failed to connect to raw client: %v", err)
	}

	address := common.HexToAddress(addressStr)
	balance, err := client.BalanceAt(ctx, address, nil)
	if err != nil {
		log.Fatalf("failed to get balance of address %s: %v", addressStr, err)
	}

	fmt.Printf("BalanceAtByAddress: %d(wei)\n", balance.Int64())
	fmt.Printf("BalanceAtByAddress: %s(ETH)\n", new(big.Rat).SetFrac(balance, big.NewInt(1e18)).FloatString(4))

}
