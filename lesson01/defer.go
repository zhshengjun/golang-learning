package main

import "fmt"

func main() {
	// 就是将 defer 的逻辑压入栈，栈遵循「后进先出」的逻辑
	defer fmt.Println("这是第一行") // 最后打印
	defer fmt.Println("这是第二行")
	defer fmt.Println("这是第三行") // 最先打印

	fmt.Println("这是 正常的函数逻辑")
}
