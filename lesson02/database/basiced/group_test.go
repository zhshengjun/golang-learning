package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestGroup(t *testing.T) {
	common.ExecSql(func(db *gorm.DB) {
		var counts []models.StatusCount
		db.Model(&models.User{}).Select("status", "count(*) as count").
			Group("status").Scan(&counts)

		common.FormatePrint(counts)
	})
}
