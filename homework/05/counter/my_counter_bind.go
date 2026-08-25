package counter

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// IncreaseByBind 使用 bind.BoundContract 调用 increase()。
func IncreaseByBind(privateKeyHex string, contractHex string) (*big.Int, error) {
	rpcURL := os.Getenv("RAW_URL")
	if rpcURL == "" {
		return nil, fmt.Errorf("RAW_URL is empty")
	}
	if !common.IsHexAddress(contractHex) {
		return nil, fmt.Errorf("invalid contract address: %s", contractHex)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("connect RPC: %w", err)
	}
	defer client.Close()

	privateKeyHex = strings.TrimPrefix(strings.TrimSpace(privateKeyHex), "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("create transact opts: %w", err)
	}
	auth.Context = ctx

	parsedABI, err := abi.JSON(strings.NewReader(CounterJson))
	if err != nil {
		return nil, fmt.Errorf("parse contract ABI: %w", err)
	}

	contract := common.HexToAddress(contractHex)
	bound := bind.NewBoundContract(contract, parsedABI, client, client, client)
	tx, err := bound.Transact(auth, "increase")
	if err != nil {
		return nil, fmt.Errorf("send increase transaction: %w", err)
	}

	waitCtx, cancelWait := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelWait()
	receipt, err := bind.WaitMined(waitCtx, client, tx)
	if err != nil {
		return nil, fmt.Errorf("wait for receipt: %w", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("transaction failed, status=%d", receipt.Status)
	}

	value, err := parseIncreasedValue(receipt, contract, auth.From)
	if err != nil {
		return nil, fmt.Errorf("parse Increased event: %w", err)
	}
	return value, nil
}
