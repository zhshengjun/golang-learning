package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestUpdates(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {

		query := models.User{ID: 17, Age: 33, Status: true}

		db.Updates(&query)
	})
}
