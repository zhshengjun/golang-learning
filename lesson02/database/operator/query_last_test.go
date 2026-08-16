package operator

import (
	"database/utils"
	"testing"
)

func TestQueryLast(t *testing.T) {

	db := utils.InitDB()
	var user utils.User
	db.Where(utils.User{Name: "张三"}).
		Last(&user)

	utils.FormatePrint(user)
}
