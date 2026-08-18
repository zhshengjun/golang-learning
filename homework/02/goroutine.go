package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	wg.Go(func() {
		for i := 1; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Printf("协程1打印偶数：%d\n", i)
			}
		}
	})

	wg.Go(func() {
		for i := 1; i <= 10; i++ {
			if i%2 != 0 {
				fmt.Printf("协程2打印奇数：%d\n", i)
			}
		}
	})

	wg.Wait()
}
