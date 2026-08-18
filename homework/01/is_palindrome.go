package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(isPalindrome2(12321))
}

func isPalindrome(num int) bool {
	if num < 0 {
		return false
	}
	var digits []int
	for _, c := range strconv.Itoa(num) {
		digits = append(digits, int(c-'0'))
	}
	left := 0
	right := len(digits) - 1

	for left <= right {
		if digits[left] != digits[right] {
			return false
		}
		left++
		right--
	}
	return true
}

func isPalindrome2(num int) bool {
	if num < 0 {
		return false
	}
	temp := 0
	result := num
	for result > 0 {
		temp = temp*10 + result%10
		result /= 10
	}
	return num == temp
}
