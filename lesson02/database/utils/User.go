package utils

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"  json:"id"`
	Name      string    `gorm:"size:255"  json:"name"`
	Age       int       `gorm:"default:0"  json:"age"`
	Status    bool      `gorm:"default:false"  json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime"  json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"  json:"updated_at"`
}

type StatusCount struct {
	Status bool  `json:"status" gorm:"column:status"`
	Count  int64 `json:"count" gorm:"column:count"`
}
