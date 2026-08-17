package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestUpdateOne(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {

		query := models.User{ID: 17}

		db.Model(&query).Update("age", 44)
	})
}
