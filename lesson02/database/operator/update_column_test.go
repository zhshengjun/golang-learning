package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestUpdateColumn(t *testing.T) {

	utils.Exec(func(db *gorm.DB) {

		query := utils.User{ID: 18}

		db.Model(&query).UpdateColumn("age", 66)
	})
}
