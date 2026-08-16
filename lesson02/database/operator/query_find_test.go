package operator

import (
	"database/utils"
	"testing"
)

func TestQueryFind(t *testing.T) {

	db := utils.InitDB()
	var users []utils.User
	db.
		//Where(utils.User{Name: "张三"}).
		Find(&users)

	utils.FormatePrint(users)
}
