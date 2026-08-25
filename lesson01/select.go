package main

import (
	"fmt"
	"time"
)

func blockChannel() {
	ch := make(chan int, 1)
	ch2 := make(chan int, 1)

	go func() {
		ch <- 1
	}()

	go func() {
		ch2 <- 2
	}()
	time.Sleep(1 * time.Second)

	select {
	case value := <-ch:
		fmt.Println(value)
	case value1 := <-ch2:
		fmt.Println(value1)
	}

}

func nonBlockChannel() {
	ch := make(chan int)

	select {
	case ch <- 42:
		fmt.Println("Sent value")
	case <-time.After(2 * time.Second):
		fmt.Println("timed out")
	default:
		fmt.Println("channel is full")
	}

	select {
	case value := <-ch:
		fmt.Println(value)
	default:
		fmt.Println("No value received")
	}

}
func channelBlock() {
	ch := make(chan int)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println(<-ch) // 接收者
	}()

	ch <- 42
	time.Sleep(2 * time.Second)
}

func nonBlockChannel2() {
	ch := make(chan int)
	defer close(ch)

	select {
	case value := <-ch:
		fmt.Println(value)
	}

	time.Sleep(2 * time.Second)

	select {
	case ch <- 42:
		fmt.Println("Sent value")
	case <-time.After(2 * time.Second):
		fmt.Println("timed out")
	}

}

func nonBlockChannel3(number int) {
	ch := make(chan int, 3)
	ch2 := make(chan int, 3)

	go func() {
		defer close(ch)
		for i := 0; i < number; i++ {
			//time.Sleep(1 * time.Second)
			ch <- i
		}

	}()

	go func() {
		defer close(ch2)
		for i := 0; i < number; i++ {
			//time.Sleep(1 * time.Second)
			ch2 <- i
		}

	}()

	for {
		select {
		case value, ok := <-ch:
			if !ok {
				ch = nil
				continue
			}
			fmt.Println(value)
		case value, ok := <-ch2:
			if !ok {
				ch2 = nil
				continue
			}
			fmt.Println(value)
		default:
			if ch == nil || ch2 == nil {
				fmt.Println("No value received")
				return
			}

		}
	}

}

func nonBlockChannel4(number int) {
	ch := make(chan int)
	ch2 := make(chan int)

	go func() {
		defer close(ch)
		for i := 0; i <= number; i++ {
			//time.Sleep(1 * time.Second)
			ch <- i
		}

	}()

	go func() {
		defer close(ch2)
		for i := 0; i <= number; i++ {
			//time.Sleep(1 * time.Second)
			ch2 <- i
		}

	}()

	for ch != nil || ch2 != nil {
		select {
		case value, ok := <-ch:
			if !ok {
				ch = nil
				continue
			}
			fmt.Println(value)
		case value, ok := <-ch2:
			if !ok {
				ch2 = nil
				continue
			}
			fmt.Println(value)
		}
	}

	fmt.Println("程序结束")

}

func main() {

	//blockChannel()
	//nonBlockChannel()
	//channelBlock()
	//nonBlockChannel3(20)
	nonBlockChannel4(20)
}
