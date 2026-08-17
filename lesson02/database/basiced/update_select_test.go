package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestUpdateSelect(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {

		query := models.User{ID: 20, Name: "test", Age: 51, Status: true}

		db.Model(&query).Select("age", "status").Updates(&query)
	})
}
