package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestGroup(t *testing.T) {
	utils.ExecSql(func(db *gorm.DB) {
		var counts []utils.StatusCount
		db.Model(&utils.User{}).Select("status", "count(*) as count").
			Group("status").Scan(&counts)

		utils.FormatePrint(counts)
	})
}
