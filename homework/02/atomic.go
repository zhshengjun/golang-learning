package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var counter atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Go(func() {
			counter.Add(1)
		})
	}
	wg.Wait()
	fmt.Println("Counter:", counter.Load())
}
