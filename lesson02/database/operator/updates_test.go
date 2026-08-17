package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestUpdates(t *testing.T) {

	utils.Exec(func(db *gorm.DB) {

		query := utils.User{ID: 17, Age: 33, Status: true}

		db.Updates(&query)
	})
}
