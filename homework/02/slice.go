package main

import "fmt"

func main() {
	nums := []int{3, 2, 3}
	multiplyTwo(nums)
	fmt.Println(nums)
}

func multiplyTwo(nums []int) {
	for i := range nums {
		nums[i] *= 2
	}
}
