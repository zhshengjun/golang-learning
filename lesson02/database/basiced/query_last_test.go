package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestQueryLast(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {
		var user models.User
		db.Where(models.User{Name: "张三"}).
			Last(&user)

		common.FormatePrint(user)
	})
}
