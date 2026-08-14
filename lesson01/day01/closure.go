package main

import "fmt"

func counter() func() int {
	count := 0

	return func() int {
		count++
		return count
	}
}

func multiplier(multiplier int) func(int) int {
	return func(base int) int {
		return multiplier * base
	}
}

func multiplier2(one int) func(int) func(int) int {
	return func(two int) func(int) int {
		return func(three int) int {
			return one * two * three
		}
	}
}

func wallet(initBalance int) (deposit, withdraw func(int), balance func() int) {

	amountBalance := initBalance

	deposit = func(amount int) {
		amountBalance += amount
	}
	withdraw = func(amount int) {
		amountBalance -= amount
	}
	balance = func() int {
		return amountBalance
	}
	return
}

func main() {
	c := counter()

	fmt.Println(c()) // 1
	fmt.Println(c()) // 2
	fmt.Println(c()) // 3

	fmt.Println("==================")

	double := multiplier(2)
	fmt.Println(double(2))

	triple := multiplier(3)
	fmt.Println(triple(4))
	fmt.Println(multiplier(3)(5))

	fmt.Println("==================")
	one := multiplier2(3)

	fmt.Printf("one的类型 %T\n", one)
	two := one(4)

	fmt.Printf("two %T\n", two)
	i := two(5)

	fmt.Printf("i的类型 %T\n", i)
	fmt.Println(i)
	fmt.Println(two(5))

	fmt.Println(two(6))

	fmt.Println(multiplier2(3)(4)(5))

	fmt.Println("========金融账户==========")

	deposit, withdraw, balance := wallet(100)
	deposit(100)
	fmt.Println(balance())

	withdraw(60)
	fmt.Println(balance())

}
