package main

import "fmt"

func main() {
	target := 5
	fmt.Println(target)
	sum(&target)
	fmt.Println(target)

}

func sum(target *int) int {
	*target = *target + 10
	return *target
}
