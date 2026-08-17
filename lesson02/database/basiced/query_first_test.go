package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestQueryFirst2(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {
		var user models.User
		db.Where(models.User{Name: "张三"}).
			First(&user)

		common.FormatePrint(user)
	})
}
