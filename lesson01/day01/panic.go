package main

import "fmt"

func main() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("这是 recover 捕获的错误：%s\n", err)
		}
	}()

	fmt.Println("这是正常逻辑，panic 之前")
	panic("这是逻辑主动抛出panic错误")
	//fmt.Println("panic 之后")
}
