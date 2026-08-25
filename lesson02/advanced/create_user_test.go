package advanced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {

		users := []models.User{
			{Name: "张三", Age: 18, Status: true},
			{Name: "李四", Age: 19},
			{Name: "王二", Age: 20},
			{Name: "麻子", Age: 21, Status: false},
		}
		db.Create(&users)
	})

}
