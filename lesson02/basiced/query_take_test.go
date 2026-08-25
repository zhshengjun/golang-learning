package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestQueryTake(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {
		var user models.User
		db.Where(models.User{Name: "张三"}).
			Take(&user)

		common.FormatePrint(user)
	})
}
