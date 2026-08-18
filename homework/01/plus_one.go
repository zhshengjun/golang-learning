package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	fmt.Println(plusOne([]int{1, 2, 3, 9, 9, 9}))
}

func plusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}

	return append([]int{1}, digits...)
}

func plusOne2(digits []int) []int {
	n := len(digits)
	var ans float64
	for i := 0; i < n; i++ {
		ans = ans + float64(digits[i])*math.Pow(10, float64(n-i-1))
	}
	ans = ans + 1

	s := strconv.Itoa(int(ans))
	digits = make([]int, len(s))

	for i, c := range s {
		digits[i] = int(c - '0')
	}

	return digits
}
