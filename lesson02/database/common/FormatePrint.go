package common

import (
	"encoding/json"
	"fmt"
)

func FormatePrint(params any) {
	data, _ := json.MarshalIndent(params, "", "  ")
	fmt.Printf("当前操作的结果：%s\n", string(data))
}
