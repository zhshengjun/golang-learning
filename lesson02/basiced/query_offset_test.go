package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestQueryOffset(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {
		var queryList []models.User
		db.Where(models.User{Age: 20}).Offset(0).Limit(2).Find(&queryList)
		common.FormatePrint(queryList)
	})
}
