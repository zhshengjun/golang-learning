package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestMigrate(t *testing.T) {
	common.ExecSql(func(db *gorm.DB) {
		err := db.AutoMigrate(&models.User{})
		if err != nil {
			return
		}
	})
}
