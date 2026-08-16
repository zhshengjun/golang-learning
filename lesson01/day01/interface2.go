package main

import (
	"fmt"
)

type WalletInterface2 interface {
	Deposit(amount int)
	Withdraw(amount int)
}

type Wallet2 struct {
	balance int
}

func (wallet Wallet2) Deposit(amount int) {
	wallet.balance += amount
}
func (wallet Wallet2) Withdraw(amount int) {
	if wallet.balance < amount {
		panic("余额不足")
	}
	wallet.balance -= amount
}

func main() {

	wallet := Wallet2{balance: 1000}
	wallet.Deposit(300)
	wallet.Withdraw(100)

	fmt.Printf("类型：%T\n", wallet)

	var x WalletInterface2 = wallet
	fmt.Println(x)
	var y WalletInterface2 = &wallet
	fmt.Println(y)
}
