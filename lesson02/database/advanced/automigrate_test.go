package advanced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestMigrateAssociation(t *testing.T) {
	common.ExecSql(func(db *gorm.DB) {

		if err := db.SetupJoinTable(
			&models.User{},
			"Roles",
			&models.UserRole{},
		); err != nil {
			panic(err)
		}

		if err := db.SetupJoinTable(
			&models.Role{},
			"Users",
			&models.UserRole{},
		); err != nil {
			panic(err)
		}

		if err := db.AutoMigrate(
			models.User{},
			&models.Profile{},
			&models.Product{},
			&models.Order{},
			&models.OrderItem{},
			&models.Role{},
		); err != nil {
			panic(err)
		}
	})
}
