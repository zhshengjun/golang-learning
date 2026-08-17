package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestUpdateColumns(t *testing.T) {

	utils.Exec(func(db *gorm.DB) {

		query := utils.User{ID: 19}

		db.Model(&query).UpdateColumns(map[string]any{
			"age":    89,
			"status": true,
		})
	})
}
