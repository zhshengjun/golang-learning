package main

import (
	"fmt"
	"sync"
	"time"
)

func channel() {
	ch := make(chan int)
	go func() { ch <- 0 }()

	fmt.Println(<-ch)
}

func channelBuffer() {
	ch := make(chan int, 3)
	go func() { ch <- 0 }()
	go func() { ch <- 1 }()
	go func() { ch <- 2 }()

	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

func channelClose() {
	ch := make(chan int, 3)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		ch <- 0
		defer wg.Done()
	}()
	go func() {
		ch <- 1
		defer wg.Done()
	}()
	go func() {
		ch <- 2
		defer wg.Done()
	}()
	wg.Wait()
	// 关闭后只能读取，不能写入了
	close(ch)

	for n := range ch {
		fmt.Println(n)
	}

}

func channelCloseBuffer() {
	fmt.Println("函数开始")
	ch := make(chan int, 2)
	defer close(ch)

	go func() {
		for n := range ch {
			fmt.Printf("协程1读取：%v\n", n)
			time.Sleep(1 * time.Second)
		}

	}()

	go func() {
		for n := range ch {
			fmt.Printf("协程2读取：%v\n", n)
			time.Sleep(1 * time.Second)
		}

	}()

	for i := 0; i < 30; i++ {
		ch <- i
	}
	time.Sleep(15 * time.Second)
	fmt.Printf("函数结束")
}

func main() {

	//channel()
	//channelBuffer()
	//channelClose()
	channelCloseBuffer()
}
