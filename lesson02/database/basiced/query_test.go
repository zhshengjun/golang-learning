package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestQuery(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {

		var query models.User
		db.Where(models.User{Age: 18}).First(&query)
		common.FormatePrint(query)
	})
}
