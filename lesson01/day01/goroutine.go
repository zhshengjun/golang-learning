package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}

	for i := 11; i <= 20; i++ {
		wg.Go(func() {
			defer wg.Done()
			fmt.Println(i)
		})
	}

	wg.Wait()
}
