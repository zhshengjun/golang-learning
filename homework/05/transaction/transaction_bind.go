package transaction

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

// TransferTokenByBind 使用 bind.BoundContract 调用 ERC-20 transfer。
func TransferTokenByBind(privateKeyHex string, toAddressHex string, contractHex string, amount float64) error {
	rpcURL := os.Getenv("RAW_URL")
	if rpcURL == "" {
		return fmt.Errorf("RAW_URL is empty")
	}
	if !common.IsHexAddress(toAddressHex) {
		return fmt.Errorf("invalid recipient address: %s", toAddressHex)
	}
	if !common.IsHexAddress(contractHex) {
		return fmt.Errorf("invalid token contract address: %s", contractHex)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return fmt.Errorf("connect RPC: %w", err)
	}
	defer client.Close()

	privateKeyHex = strings.TrimPrefix(strings.TrimSpace(privateKeyHex), "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("parse private key: %w", err)
	}
	from := crypto.PubkeyToAddress(privateKey.PublicKey)
	to := common.HexToAddress(toAddressHex)
	contract := common.HexToAddress(contractHex)

	amountWei, _ := new(big.Float).
		Mul(big.NewFloat(amount), big.NewFloat(1e18)).
		Int(new(big.Int))

	parsedABI, err := abi.JSON(strings.NewReader(Erc20Json))
	if err != nil {
		return fmt.Errorf("parse ERC-20 ABI: %w", err)
	}
	bound := bind.NewBoundContract(contract, parsedABI, client, client, client)

	// balanceOf 是只读方法，所以使用 Call，而不是 Transact。
	var outputs []interface{}
	if err := bound.Call(
		&bind.CallOpts{Context: ctx, From: from},
		&outputs,
		"balanceOf",
		from,
	); err != nil {
		return fmt.Errorf("call balanceOf: %w", err)
	}
	if len(outputs) != 1 {
		return fmt.Errorf("invalid balanceOf result: got %d values", len(outputs))
	}
	tokenBalance, ok := outputs[0].(*big.Int)
	if !ok {
		return fmt.Errorf("unexpected balanceOf result type: %T", outputs[0])
	}
	if amountWei.Cmp(tokenBalance) > 0 {
		return fmt.Errorf("insufficient token balance: balance=%s, amount=%s", tokenBalance, amountWei)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("get chain ID: %w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("create transact opts: %w", err)
	}
	auth.Context = ctx

	tx, err := bound.Transact(auth, "transfer", to, amountWei)
	if err != nil {
		return fmt.Errorf("send transfer transaction: %w", err)
	}

	waitCtx, cancelWait := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelWait()
	receipt, err := bind.WaitMined(waitCtx, client, tx)
	if err != nil {
		return fmt.Errorf("wait for receipt: %w", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("transaction failed, status=%d", receipt.Status)
	}

	fmt.Printf("bind txHash: %s\n", tx.Hash().Hex())
	fmt.Printf("bind Tx Status: %d\n", receipt.Status)
	fmt.Printf("bind Tx GasUsed: %d\n", receipt.GasUsed)
	return nil
}
