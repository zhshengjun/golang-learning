package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestQueryOffset(t *testing.T) {

	utils.Exec(func(db *gorm.DB) {
		var queryList []utils.User
		db.Where(utils.User{Age: 20}).Offset(0).Limit(2).Find(&queryList)
		utils.FormatePrint(queryList)
	})
}
