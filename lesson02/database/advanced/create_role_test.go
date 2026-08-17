package advanced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestCreateRole(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {

		roles := []models.Role{
			{RoleName: "角色01"},
			{RoleName: "角色02"},
			{RoleName: "角色03"},
			{RoleName: "角色04"},
		}
		db.Create(&roles)
	})

}
