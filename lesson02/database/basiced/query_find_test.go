package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestQueryFind(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {
		var users []models.User
		db.
			//Where(utils.User{Name: "张三"}).
			Find(&users)

		common.FormatePrint(users)
	})
}
