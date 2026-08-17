package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestQueryFindInBatch(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {
		var users []models.User
		db.
			Where(models.User{Name: "张三"}).
			FindInBatches(&users, 1, func(tx *gorm.DB, batch int) error {
				return nil
			})

		common.FormatePrint(users)
	})
}
