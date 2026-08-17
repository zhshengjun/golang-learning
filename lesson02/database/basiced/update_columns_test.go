package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestUpdateColumns(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {

		query := models.User{ID: 19}

		db.Model(&query).UpdateColumns(map[string]any{
			"age":    89,
			"status": true,
		})
	})
}
