package main

import (
	"fmt"
	"time"
)

func main() {

	channel()
	channel2()

}

func channel() {
	ch := make(chan int)
	done := make(chan bool)

	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
			if i == 9 {
				done <- true
			}
		}
	}()

	for {
		select {
		case i := <-ch:
			fmt.Printf("方法打印：%d\n", i)
		case <-done:
			fmt.Println("done")
			return
		}
	}
}

func channel2() {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for i := range ch {
		fmt.Printf("方法2打印：%d\n", i)
	}

	time.Sleep(time.Second)
}
