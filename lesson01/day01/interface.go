package main

import (
	"fmt"
	"math"
)

type WalletInterface interface {
	Deposit(amount int)

	Withdraw(amount int)

	Balance() int
}

type Wallet struct {
	balance int
}

func (wallet *Wallet) Deposit(amount int) {
	wallet.balance += amount
}
func (wallet *Wallet) Withdraw(amount int) {
	if wallet.balance < amount {
		panic("余额不足")
	}
	wallet.balance -= amount
}

func (wallet *Wallet) Balance() int {
	return wallet.balance
}

type Shape interface {
	Area() float64
}

type Circle struct {
	radius float64
}

type Rectangle struct {
	length, width float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (r Rectangle) Area() float64 {
	return r.length * r.width
}

func main() {

	wallet := Wallet{balance: 1000}
	wallet.Deposit(300)
	wallet.Withdraw(100)
	fmt.Println(wallet.Balance())

	fmt.Printf("类型：%T\n", wallet)

	wallet1 := &Wallet{balance: 2000}
	wallet1.Deposit(300)
	wallet1.Withdraw(100)
	fmt.Println(wallet1.Balance())

	fmt.Printf("类型：%T\n", wallet1)

	wallet2 := Wallet{balance: 3000}
	(&wallet2).Deposit(300)
	(&wallet2).Withdraw(100)
	fmt.Println((&wallet2).Balance())

	var wallet3 = Wallet{balance: 4000}
	var wallet4 = &wallet3
	wallet4.Deposit(300)
	wallet4.Withdraw(100)
	fmt.Println(wallet3.Balance())

	fmt.Println("=======接口======")

	shapes := []Shape{&Circle{radius: 5}, &Rectangle{length: 4, width: 5}}
	for _, shape := range shapes {
		fmt.Println(shape.Area())
	}

	var x WalletInterface = &wallet3
	fmt.Println(x)

	var c Circle = Circle{radius: 5}

	var y Shape = &c
	fmt.Println(y)
}
