package transaction

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// TransferToken 操作指定的token
func TransferToken(privateKeyHex string, toAddressHex string, contractHex string, amount float64) {

	//1. 连接
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()

	client, _ := ethclient.DialContext(ctx, os.Getenv("RAW_URL"))
	defer client.Close()

	//2. 私钥生成公钥key，用于余额获取
	privateKey, _ := crypto.HexToECDSA(privateKeyHex)
	publicKey := privateKey.PublicKey
	from := crypto.PubkeyToAddress(publicKey)

	// 转账目标地址
	to := common.HexToAddress(toAddressHex)
	// 合约地址
	contract := common.HexToAddress(contractHex)

	//2. nonce + chainId
	nonce, _ := client.PendingNonceAt(ctx, from)
	chainID, _ := client.ChainID(ctx)

	//3. bassFee 注意区分EIP-1559，gasTipCap,gasFeeCap
	header, _ := client.HeaderByNumber(ctx, big.NewInt(int64(rpc.LatestBlockNumber)))

	baseFee := header.BaseFee
	if baseFee == nil {
		suggestGasPrice, _ := client.SuggestGasPrice(ctx)
		baseFee = suggestGasPrice
	}
	gasTipCap, _ := client.SuggestGasTipCap(ctx)

	gasFeeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), gasTipCap)

	amountWei, _ := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e18)).Int(new(big.Int))
	// 构建 token abi
	parsedABI, _ := abi.JSON(strings.NewReader(Erc20Json))

	//4. gasLimit
	contractData, _ := parsedABI.Pack("transfer", to, amountWei)
	gasLimit, _ := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    &contract, // 注意这里是合约地址
		Value: big.NewInt(0),
		Data:  contractData,
	})

	// 这里可以适当增加 gas，
	gasLimit = gasLimit * 120 / 100

	maxGasPrice := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasFeeCap)

	//5. 校验gas余额 和 token余额
	balance, err := client.BalanceAt(ctx, from, nil)
	if err != nil {
		log.Fatalf("Error getting balance: %v", err)
	}
	if balance.Cmp(maxGasPrice) == -1 {
		log.Fatalf("balance is not enough for apy gas!")
	}

	tokenBalanceData, _ := parsedABI.Pack("balanceOf", from)

	callContract, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &contract,
		Data: tokenBalanceData,
	}, nil)
	if err != nil {
		log.Fatalf("Error getting contract balance: %v", err)
	}

	tokenValues, _ := parsedABI.Unpack("balanceOf", callContract)
	if len(tokenValues) != 1 {
		log.Fatalf("token amount is not enough!")
	}
	if amountWei.Cmp(tokenValues[0].(*big.Int)) > 0 {
		log.Fatalf("amount is not enough!")
	}

	//6. 构建动态交易体 txData 和 tx
	dynamicFeeTx := types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &contract,
		Value:     big.NewInt(0),
		Data:      contractData,
	}
	tx := types.NewTx(&dynamicFeeTx)

	//7. 获取签名
	transaction, _ := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)

	//8. 发送交易
	err = client.SendTransaction(ctx, transaction)
	if err != nil {
		log.Fatal("failed to send transaction")
	}
	txHash := transaction.Hash()
	fmt.Println("txHash:", txHash.Hex())
	//9. 通过hash查询回执

	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err != nil {
			fmt.Printf("failed to get transaction receipt: %v\n", err)
			time.Sleep(3 * time.Second)
			continue
		}
		fmt.Printf("Tx Status: %d\n", receipt.Status)
		fmt.Printf("Tx GasUsed: %d\n", receipt.GasUsed)
		return
	}

}
