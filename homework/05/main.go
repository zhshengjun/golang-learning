package main

import (
	"fmt"
	"homework05/counter"
	"homework05/transaction"
	"log"
	"os"
)

func main() {

	// 转账合约，这里没有自己部署新的代币，使用了 LINK，合约地址不一样而已
	// 这里不使用 bind.Transact 手搓
	transaction.TransferToken(os.Getenv("ACCOUNT_FROM_PRIVATE"), os.Getenv("ACCOUNT_TO_ADDRESS"), "0x779877A7B0D9E8603169DdbD7836e478b4624789", 1)
	_ = transaction.TransferTokenByBind(os.Getenv("ACCOUNT_FROM_PRIVATE"), os.Getenv("ACCOUNT_TO_ADDRESS"), "0x779877A7B0D9E8603169DdbD7836e478b4624789", 2)

	// 这里使用 bind.Transact

	//counter.DeployCounterOperator()
	//contractAddress := "0x6357571E8dAD3901eF6201402fD2e95AdeA7F0B6"
	//IncreaseByBind(os.Getenv("ACCOUNT_FROM_PRIVATE"), contractAddress)
	//currentCounter(contractAddress)

}

func currentCounter(contractAddress string) {
	current, err := counter.GetCounter(contractAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current count: %s\n", current.String())
}

func increased(privateHex string, contractAddress string) {
	increased, err := counter.Increase(privateHex, contractAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Increase count: %s\n", increased.String())
}

func IncreaseByBind(privateHex string, contractAddress string) {
	increased, err := counter.IncreaseByBind(privateHex, contractAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Increase count: %s\n", increased.String())
}
