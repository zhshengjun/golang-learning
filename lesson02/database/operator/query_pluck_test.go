package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestQueryPluck(t *testing.T) {

	utils.Exec(func(db *gorm.DB) {
		var names []string
		db.Model(&utils.User{}).Where(utils.User{Age: 20}).
			Pluck("name", &names)

		utils.FormatePrint(names)
	})
}
