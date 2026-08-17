package advanced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestQueryUser(t *testing.T) {
	common.ExecSql(func(db *gorm.DB) {
		user := models.User{}
		db.Model(&models.User{}).Where("id = ?", 30).First(&user)
		common.FormatePrint(&user)
	},
	)
}
