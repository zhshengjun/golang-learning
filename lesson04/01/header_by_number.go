package main

import (
	"01/commons"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main1() {

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
	headerNumber := big.NewInt(11527733)

	header, err := client.HeaderByNumber(ctx, headerNumber)
	if err != nil {
		log.Fatalf("failed get header: %v", err)
	}

	fmt.Printf("Chain ID:  %d\n", chainID)
	fmt.Printf("Block Number: %d\n", header.Number)
	fmt.Printf("Block Hash: %s\n", header.Hash().Hex())
	fmt.Printf("Parent Hash: %s\n", header.ParentHash.Hex())
	fmt.Printf("State Root: %s\n", header.Root.Hex())
	fmt.Printf("Block Time: %s\n", time.Unix(int64(header.Time), 0).Format(time.RFC3339))
}

// 下面这段代码是为了演示， block的 hash可能不同，因为header中没有 Hash 属性，只有 Hash() 函数，获取是实时计算，可能受环境影响，可能和页面展示的不一致。目前我环境的是和 RPC 返回的一致
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
	num := big.NewInt(11527673)
	_ = new(big.Int).Sub(num, big.NewInt(0))
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("failed get header: %v", err)
	}

	fmt.Printf("RPC URL: %s\n", commons.Endpoint)
	fmt.Printf("Chain ID:  %d\n", chainID)
	fmt.Printf("Block Number: %d\n", header.Number)
	fmt.Printf("Block Hash: %s\n", header.Hash().Hex())
	fmt.Printf("Parent Hash: %s\n", header.ParentHash.Hex())
	fmt.Printf("Block Time: %s\n", time.Unix(int64(header.Time), 0).Format(time.RFC3339))
	fmt.Println("+++++++++++++++++++++++++++")

	// 这里获取指定的bloc
	//num := new(big.Int).Sub(header.Number, big.NewInt(0))
	//header, err = client.HeaderByNumber(ctx, num)
	//fmt.Printf("Header Number: %d\n", header.Number)

	// c查询 safe 块信息

	safeHeader, safeHash, err := getBlockByTag(ctx, client, "safe")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Safe Block Number : %d\n", safeHeader.Number.Uint64())
	fmt.Printf("Safe Block Hash（RPC提供的） : %s\n", safeHash.Hex())
	fmt.Printf("Safe Block Hash (计算出来的) : %s\n", safeHeader.Hash().Hex())
	fmt.Printf("Safe Block Time: %s\n", time.Unix(int64(safeHeader.Time), 0).Format(time.RFC3339))
	fmt.Printf("Safe with Latest Confirmations Number: %d\n", header.Number.Int64()-safeHeader.Number.Int64())

	finalizedHeader, finalizedHash, err := getBlockByTag(ctx, client, "finalized")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Finalized Block Number : %d\n", finalizedHeader.Number.Uint64())
	fmt.Printf("Finalized Block Hash（RPC提供的） : %s\n", finalizedHash.Hex())
	fmt.Printf("Finalized Block Hash (计算出来的) : %s\n", finalizedHeader.Hash().Hex())
	fmt.Printf("Finalized Block Time: %s\n", time.Unix(int64(finalizedHeader.Time), 0).Format(time.RFC3339))
	fmt.Printf("Finalized By Lates Confirmations Number: %d\n", header.Number.Int64()-finalizedHeader.Number.Int64())

}

func getBlockByTag(ctx context.Context, client *ethclient.Client, tag string) (*types.Header, common.Hash, error) {

	rpcClient := client.Client()

	var raw json.RawMessage

	err := rpcClient.CallContext(ctx, &raw, "eth_getBlockByNumber", tag, false)

	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("rpc call error %w", err)
	}

	if len(raw) == 0 || string(raw) == "null" {
		return nil, common.Hash{}, fmt.Errorf("rpc call error %w", errors.New("rpc call error null"))
	}

	fmt.Println(string(raw))

	var blockData struct {
		Number      string         `json:"number"`
		Hash        common.Hash    `json:"hash"`
		ParentHash  common.Hash    `json:"parentHash"`
		UncleHash   common.Hash    `json:"sha3Uncles"`
		Coinbase    common.Address `json:"miner"`
		Root        common.Hash    `json:"stateRoot"`
		TxHash      common.Hash    `json:"transactionsRoot"`
		ReceiptHash common.Hash    `json:"receiptRoot"`
		LogsBloom   hexutil.Bytes  `json:"logsBloom"`
		Difficulty  *hexutil.Big   `json:"difficulty"`
		GasLimit    hexutil.Uint64 `json:"gasLimit"`
		GasUsed     hexutil.Uint64 `json:"gasUsed"`
		Time        hexutil.Uint64 `json:"timestamp"`
		ExtraData   hexutil.Bytes  `json:"extraData"`
		MixDigest   common.Hash    `json:"mixHash"`
		Nonce       hexutil.Bytes  `json:"nonce"`
		BaseFee     *hexutil.Big   `json:"baseFee"`
	}
	if err = json.Unmarshal(raw, &blockData); err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to unmarshal header: %w", err)
	}

	num, ok := new(big.Int).SetString(blockData.Number[2:], 16)
	if !ok {
		return nil, common.Hash{}, fmt.Errorf("failed to convert number to big.Int %s", blockData.Number[2:])
	}

	// 构造完整的Header
	header := &types.Header{
		ParentHash:  blockData.ParentHash,
		UncleHash:   blockData.UncleHash,
		Coinbase:    blockData.Coinbase,
		Root:        blockData.Root,
		TxHash:      blockData.TxHash,
		ReceiptHash: blockData.ReceiptHash,
		Bloom:       types.BytesToBloom(blockData.LogsBloom),
		Difficulty:  big.NewInt(0),
		Number:      num,
		GasLimit:    uint64(blockData.GasLimit),
		GasUsed:     uint64(blockData.GasUsed),
		Time:        uint64(blockData.Time),
		Extra:       blockData.ExtraData,
		MixDigest:   blockData.MixDigest,
		BaseFee:     nil,
	}

	// set Difficulty
	if blockData.Difficulty != nil {
		header.Difficulty = blockData.Difficulty.ToInt()
	}

	//set BaseFee
	if blockData.BaseFee != nil {
		header.BaseFee = blockData.BaseFee.ToInt()
	}

	// set nonce
	if len(blockData.Nonce) >= 8 {
		var nonceBytes [8]byte
		copy(nonceBytes[:], blockData.Nonce[:8])
		header.Nonce = nonceBytes
	}

	// 返回

	return header, blockData.Hash, nil

}
