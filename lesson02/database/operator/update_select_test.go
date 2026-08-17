package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestUpdateSelect(t *testing.T) {

	utils.ExecSql(func(db *gorm.DB) {

		query := utils.User{ID: 20, Name: "test", Age: 51, Status: true}

		db.Model(&query).Select("age", "status").Updates(&query)
	})
}
