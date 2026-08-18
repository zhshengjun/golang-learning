package main

import "fmt"

/*
*
给你一个 非空 整数数组 nums ，除了某个元素只出现一次以外，其余每个元素均出现多次。
找出那个只出现了一次的元素。
*/
func main() {
	nums := []int{1, 2, 2, 4, 4, 5, 6, 7, 7, 8, 9}
	fmt.Println(singleNumber(nums))
}

func singleNumber(nums []int) int {

	numMap := map[int]int{}
	for _, key := range nums {
		numMap[key]++
	}
	for key, num := range numMap {
		if num == 1 {
			return key
		}
	}
	return 0
}
