package advanced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestAssociations(t *testing.T) {
	common.ExecSql(func(db *gorm.DB) {
		user := models.User{ID: 1}
		db.Model(&user).First(&user)

		role := models.Role{ID: 1}

		err := db.Model(&user).Association("Roles").Append(&role)
		if err != nil {
			return
		}
	})
}
