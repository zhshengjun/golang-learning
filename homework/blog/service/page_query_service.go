package service

import "gorm.io/gorm"

func Paginate(currentPage int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((currentPage - 1) * pageSize).Limit(pageSize)
	}
}

func Pages(total int, pageSize int) int {

	if total == 0 || pageSize == 0 {
		return 0
	}

	return max(1, (total+pageSize-1)/pageSize)
}
