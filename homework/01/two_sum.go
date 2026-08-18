package main

import "fmt"

func main() {
	nums := []int{3, 2, 3}
	target := 6

	fmt.Println(twoSum2(nums, target))
}

func twoSum(nums []int, target int) []int {
	n := len(nums)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {

			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return []int{}
}

/*
*
map 记录值对应的索引
*/
func twoSum2(nums []int, target int) []int {
	res := []int{0, 0}
	mp := make(map[int]int, 8)
	for i, v := range nums {
		if val, exist := mp[target-v]; exist {
			res[0] = val
			res[1] = i
		}
		mp[v] = i
	}
	return res
}
