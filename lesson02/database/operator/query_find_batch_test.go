package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestQueryFindInBatch(t *testing.T) {

	db := utils.InitDB()
	var users []utils.User
	db.
		Where(utils.User{Name: "张三"}).
		FindInBatches(&users, 1, func(tx *gorm.DB, batch int) error {
			return nil
		})

	utils.FormatePrint(users)
}
