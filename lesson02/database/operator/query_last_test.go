package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestQueryLast(t *testing.T) {

	utils.Exec(func(db *gorm.DB) {
		var user utils.User
		db.Where(utils.User{Name: "张三"}).
			Last(&user)

		utils.FormatePrint(user)
	})
}
