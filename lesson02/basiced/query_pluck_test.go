package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestQueryPluck(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {
		var names []string
		db.Model(&models.User{}).Where(models.User{Age: 20}).
			Pluck("name", &names)

		common.FormatePrint(names)
	})
}
