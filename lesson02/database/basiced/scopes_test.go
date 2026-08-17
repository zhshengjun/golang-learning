package basiced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestScopes(t *testing.T) {
	common.ExecSql(func(db *gorm.DB) {
		var users []models.User
		db.Model(&models.User{}).
			Scopes(FilterStatus(true), adultUsers(20), Paginate(0, 3)).
			Find(&users)

		common.FormatePrint(users)
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
