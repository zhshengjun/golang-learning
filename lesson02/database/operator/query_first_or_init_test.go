package operator

import (
	"database/utils"
	"testing"
)

func TestQueryFirstOrInit(t *testing.T) {

	db := utils.InitDB()
	var user utils.User
	// attrs 是创建时的参数，注意后一次 Attrs() 会覆盖前一个
	db.Where(utils.User{Name: "张三4"}).
		Attrs(utils.User{Age: 20, Status: true}).
		//Attrs(utils.User{Status: true}).
		FirstOrCreate(&user)

	utils.FormatePrint(user)
}
