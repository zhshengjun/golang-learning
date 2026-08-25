package counter

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Increase 执行 increase()，并返回这笔交易发出的 Increased 事件中的值。
func Increase(privateKeyHex string, contractHex string) (*big.Int, error) {
	return operateContract(privateKeyHex, contractHex)
}

// GetCounter 读取最新的计数器
func GetCounter(contractHex string) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, os.Getenv("RAW_URL"))
	if err != nil {
		return nil, fmt.Errorf("connect RPC: %w", err)
	}
	defer client.Close()

	if !common.IsHexAddress(contractHex) {
		return nil, fmt.Errorf("invalid contract address: %s", contractHex)
	}

	return readCounter(ctx, client, common.HexToAddress(contractHex))
}

// operateContract 发送 increase()，等待交易确认，并解析本笔交易的事件日志。
func operateContract(privateKeyHex string, contractHex string) (*big.Int, error) {
	rpcURL := os.Getenv("RAW_URL")
	if rpcURL == "" {
		return nil, fmt.Errorf("RAW_URL is empty")
	}
	// 合约地址
	if !common.IsHexAddress(contractHex) {
		return nil, fmt.Errorf("invalid contract address: %s", contractHex)
	}

	// 初始化连接
	setupCtx, cancelSetup := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelSetup()

	client, err := ethclient.DialContext(setupCtx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("connect RPC: %w", err)
	}
	defer client.Close()

	// 解析私钥
	privateKeyHex = strings.TrimPrefix(strings.TrimSpace(privateKeyHex), "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	// 地址处理
	from := crypto.PubkeyToAddress(privateKey.PublicKey)
	contract := common.HexToAddress(contractHex)

	chainID, err := client.ChainID(setupCtx)
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}
	// 账户的 nonce
	nonce, err := client.PendingNonceAt(setupCtx, from)
	if err != nil {
		return nil, fmt.Errorf("get pending nonce: %w", err)
	}

	header, err := client.HeaderByNumber(setupCtx, big.NewInt(int64(rpc.LatestBlockNumber)))
	if err != nil {
		return nil, fmt.Errorf("get latest header: %w", err)
	}
	// 基础的 baseFee
	baseFee := header.BaseFee
	if baseFee == nil {
		baseFee, err = client.SuggestGasPrice(setupCtx)
		if err != nil {
			return nil, fmt.Errorf("get gas price: %w", err)
		}
	}

	// 每单位 Gas 的小费上限
	gasTipCap, err := client.SuggestGasTipCap(setupCtx)
	if err != nil {
		return nil, fmt.Errorf("get gas tip cap: %w", err)
	}

	// gas 的最高单价
	gasFeeCap := new(big.Int).Add(
		new(big.Int).Mul(baseFee, big.NewInt(2)),
		gasTipCap,
	)

	data, err := packContractData("increase")
	if err != nil {
		return nil, fmt.Errorf("pack increase calldata: %w", err)
	}
	// 当前操作需要的 gas 单位数量
	gasLimit, err := client.EstimateGas(setupCtx, ethereum.CallMsg{
		From:  from,
		To:    &contract,
		Value: big.NewInt(0),
		Data:  data,
	})
	if err != nil {
		return nil, fmt.Errorf("estimate gas: %w", err)
	}

	// 校验余额是否足够支付 gas 费用
	maxCost := new(big.Int).Mul(
		new(big.Int).SetUint64(gasLimit),
		gasFeeCap,
	)

	balance, err := client.BalanceAt(setupCtx, from, nil)
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}
	if balance.Cmp(maxCost) < 0 {
		return nil, fmt.Errorf(
			"insufficient balance: balance=%s, maxCost=%s",
			balance.String(),
			maxCost.String(),
		)
	}
	// 构造交易体
	txData := types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		To:        &contract,
		Value:     big.NewInt(0),
		Data:      data,
	}

	//发送交易
	tx := types.NewTx(&txData)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign transaction: %w", err)
	}

	sendCtx, cancelSend := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelSend()
	if err := client.SendTransaction(sendCtx, signedTx); err != nil {
		return nil, fmt.Errorf("send transaction: %w", err)
	}

	// 等待写入成果
	waitCtx, cancelWait := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelWait()

	// 交易回执
	receipt, err := bind.WaitMined(waitCtx, client, signedTx)
	if err != nil {
		return nil, fmt.Errorf("wait for receipt: %w", err)
	}
	// staus =1 表示成功
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("transaction failed, status=%d", receipt.Status)
	}

	// 自增log中解析出结果
	value, err := parseIncreasedValue(receipt, contract, from)
	if err != nil {
		return nil, fmt.Errorf("parse Increased event: %w", err)
	}
	return value, nil
}

// parseIncreasedValue 从日志里解析操作的值
// rawLog.Topics[0]：事件签名 Increased(address,uint256)
// rawLog.Topics[1]：caller，因为 caller 是 indexed
// rawLog.Data：value，因为 value 不是 indexed
func parseIncreasedValue(receipt *types.Receipt, contract, caller common.Address) (*big.Int, error) {
	parsedABI, err := abi.JSON(strings.NewReader(CounterJson))
	if err != nil {
		return nil, err
	}

	eventABI, ok := parsedABI.Events["Increased"]
	if !ok {
		return nil, fmt.Errorf("Increased event is missing from ABI")
	}

	for _, rawLog := range receipt.Logs {
		if rawLog == nil || rawLog.Address != contract {
			continue
		}
		if len(rawLog.Topics) == 0 || rawLog.Topics[0] != eventABI.ID {
			continue
		}
		// 当前事件的日志只有1个索引，所以 只有签名和地址
		if len(rawLog.Topics) != 2 {
			return nil, fmt.Errorf("invalid Increased event topics: got %d", len(rawLog.Topics))
		}
		// 不是调用者的地址，跳过
		if common.BytesToAddress(rawLog.Topics[1].Bytes()) != caller {
			continue
		}
		// 获取日志中的data，目前data只有一个计数值
		values, err := eventABI.Inputs.NonIndexed().Unpack(rawLog.Data)
		if err != nil {
			return nil, err
		}
		if len(values) != 1 {
			return nil, fmt.Errorf("invalid Increased event data: got %d values", len(values))
		}

		value, ok := values[0].(*big.Int)
		if !ok {
			return nil, fmt.Errorf("unexpected Increased value type")
		}
		return value, nil
	}

	return nil, fmt.Errorf("increased event not found in transaction receipt")
}

func packContractData(method string) ([]byte, error) {
	parsedABI, err := abi.JSON(strings.NewReader(CounterJson))
	if err != nil {
		return nil, err
	}
	return parsedABI.Pack(method)
}

// readCounter 读取最新的计数值
func readCounter(ctx context.Context, client *ethclient.Client, contract common.Address) (*big.Int, error) {
	parsedABI, err := abi.JSON(strings.NewReader(CounterJson))
	if err != nil {
		return nil, err
	}

	data, err := parsedABI.Pack("getCounter")
	if err != nil {
		return nil, err
	}

	output, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	values, err := parsedABI.Unpack("getCounter", output)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("getCounter returned no value")
	}

	value, ok := values[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected getCounter result type")
	}
	return value, nil
}
