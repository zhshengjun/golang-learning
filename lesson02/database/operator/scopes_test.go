package operator

import (
	"database/utils"
	"testing"

	"gorm.io/gorm"
)

func TestScopes(t *testing.T) {
	utils.ExecSql(func(db *gorm.DB) {
		var users []utils.User
		db.Model(&utils.User{}).
			Scopes(FilterStatus(true), adultUsers(20), Paginate(0, 3)).
			Find(&users)

		utils.FormatePrint(users)
	})
}

func adultUsers(age int64) func(*gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {
		return db.Where("age=?", age)
	}
}

func FilterStatus(status bool) func(*gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

func Paginate(page, pageSize int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
