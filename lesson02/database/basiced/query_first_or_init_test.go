package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestQueryFirstOrInit(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {
		var user models.User
		// attrs 是创建时的参数，注意后一次 Attrs() 会覆盖前一个
		db.Where(models.User{Name: "张三4"}).
			Attrs(models.User{Age: 20, Status: true}).
			//Attrs(utils.User{Status: true}).
			FirstOrCreate(&user)

		common.FormatePrint(user)
	})
}
