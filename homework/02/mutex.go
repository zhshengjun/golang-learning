package main

import (
	"fmt"
	"sync"
)

func main() {

	calculate()
	calculateLock()
}

func calculateLock() {
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	num := 0
	for i := 0; i < 10000; i++ {
		wg.Go(func() {
			mutex.Lock()
			num++
			defer mutex.Unlock()
		})
	}
	wg.Wait()
	fmt.Println(num)
}

func calculate() {
	wg := sync.WaitGroup{}
	num := 0
	for i := 0; i < 10000; i++ {
		wg.Go(func() {

			num++

		})
	}
	wg.Wait()
	fmt.Println(num)
}
