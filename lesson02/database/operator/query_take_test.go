package operator

import (
	"database/utils"
	"testing"
)

func TestQueryTake(t *testing.T) {

	db := utils.InitDB()
	var user utils.User
	db.Where(utils.User{Name: "张三"}).
		Take(&user)

	utils.FormatePrint(user)
}
