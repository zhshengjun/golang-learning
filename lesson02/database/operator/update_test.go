package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestUpdateOne(t *testing.T) {

	utils.ExecSql(func(db *gorm.DB) {

		query := utils.User{ID: 17}

		db.Model(&query).Update("age", 44)
	})
}
