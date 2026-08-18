package main

import (
	"fmt"
	"slices"
)

func main() {
	s := "{(())}"
	fmt.Println(validParentheses(s))
}

var left = []rune{'(', '[', '{'}
var parenthesesMap = map[rune]rune{
	')': '(',
	']': '[',
	'}': '{',
}

func validParentheses(s string) bool {
	if len(s) <= 1 {
		return false
	}

	var stack []rune
	chars := []rune(s)

	for _, parentheses := range chars {
		if slices.Contains(left, parentheses) {
			stack = append(stack, parentheses)
			continue
		}

		if len(stack) == 0 {
			return false
		}
		// 出栈
		top := stack[len(stack)-1]
		// 出栈后的栈
		stack = stack[:len(stack)-1]
		if top != parenthesesMap[parentheses] {
			return false
		}
	}
	if len(stack) > 0 {
		return false
	}
	return true
}
