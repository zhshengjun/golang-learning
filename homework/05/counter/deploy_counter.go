package counter

import (
	"context"
	"fmt"
	"homework05/counter/script"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func DeployCounterOperator() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	address, txHash, err := deployCounterOperator(
		ctx,
		os.Getenv("RAW_URL"),
		os.Getenv("ACCOUNT_FROM_PRIVATE"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("deployment transaction:", txHash.Hex())
	fmt.Println("contract address:", address.Hex())
}

func deployCounterOperator(ctx context.Context, rpcURL, privateKeyHex string) (common.Address, common.Hash, error) {
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return common.Address{}, common.Hash{}, fmt.Errorf("connect Sepolia: %w", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return common.Address{}, common.Hash{}, fmt.Errorf("get chain ID: %w", err)
	}
	if chainID.Cmp(big.NewInt(11155111)) != 0 {
		return common.Address{}, common.Hash{}, fmt.Errorf("expected Sepolia chain ID 11155111, got %s", chainID)
	}

	privateKeyHex = strings.TrimPrefix(strings.TrimSpace(privateKeyHex), "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return common.Address{}, common.Hash{}, fmt.Errorf("parse private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return common.Address{}, common.Hash{}, fmt.Errorf("create transactor: %w", err)
	}
	auth.Context = ctx

	address, tx, _, err := script.DeployMyCounter(auth, client)
	if err != nil {
		return common.Address{}, common.Hash{}, fmt.Errorf("send deployment transaction: %w", err)
	}
	if _, err = bind.WaitDeployed(ctx, client, tx); err != nil {
		return common.Address{}, common.Hash{}, fmt.Errorf("wait for deployment: %w", err)
	}

	return address, tx.Hash(), nil
}
