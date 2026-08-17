package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestQuery(t *testing.T) {

	utils.Exec(func(db *gorm.DB) {

		var query utils.User
		db.Where(utils.User{Age: 18}).First(&query)
		utils.FormatePrint(query)
	})
}
