package main

import "fmt"

func main() {
	ints := [...]string{"a", "b", "c", "d", "e"}
	fmt.Println(ints)

	for index, v := range ints {
		fmt.Println(index)
		fmt.Println(v)
	}
}
