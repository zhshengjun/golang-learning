package utils

import (
	"encoding/json"
	"fmt"
)

func FormatePrint(params any) {
	data, _ := json.MarshalIndent(params, "", "  ")
	fmt.Printf("当前的格式化结果：%s\n", string(data))
}
