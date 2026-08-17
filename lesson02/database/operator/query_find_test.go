package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestQueryFind(t *testing.T) {

	utils.Exec(func(db *gorm.DB) {
		var users []utils.User
		db.
			//Where(utils.User{Name: "张三"}).
			Find(&users)

		utils.FormatePrint(users)
	})
}
