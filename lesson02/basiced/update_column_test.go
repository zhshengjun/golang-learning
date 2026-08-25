package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestUpdateColumn(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {

		query := models.User{ID: 18}

		db.Model(&query).UpdateColumn("age", 66)
	})
}
