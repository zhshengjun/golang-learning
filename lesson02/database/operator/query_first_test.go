package operator

import (
	"database/utils"
	"testing"
)

func TestQueryFirst2(t *testing.T) {

	db := utils.InitDB()
	var user utils.User
	db.Where(utils.User{Name: "张三"}).
		First(&user)

	utils.FormatePrint(user)
}
