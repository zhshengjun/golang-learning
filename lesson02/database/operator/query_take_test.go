package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestQueryTake(t *testing.T) {

	utils.ExecSql(func(db *gorm.DB) {
		var user utils.User
		db.Where(utils.User{Name: "张三"}).
			Take(&user)

		utils.FormatePrint(user)
	})
}
