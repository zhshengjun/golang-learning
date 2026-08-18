package main

import "fmt"

func main() {
	strs := []string{"flower", "flow", "flight"}

	fmt.Println(longestCommonPrefix2(strs))
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	var result []rune
	first := strs[0]
Outer:
	for i, c := range []rune(first) {
		for index := 1; index < len(strs); index++ {
			inner := strs[index]
			if i > len(inner)-1 {
				break Outer
			}
			if []rune(inner)[i] != c {
				break Outer
			}
		}
		result = append(result, c)
	}

	return string(result)
}

func longestCommonPrefix2(strs []string) string {
	n := len(strs)
	ans := strs[0]

	for i := 1; i < n; i++ {
		j := 0
		for ; j < len(ans) && j < len(strs[i]); j++ {
			if ans[j] != strs[i][j] {
				break
			}
		}
		ans = ans[:j]
	}

	return ans
}
