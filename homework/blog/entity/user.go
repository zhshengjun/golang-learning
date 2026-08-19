package entity

import "time"

type User struct {
	ID       int64     `json:"id" gorm:"primary_key"`
	UserName string    `json:"userName" gorm:"column:user_name"`
	Password string    `json:"-"`
	Email    string    `json:"email"`
	Status   bool      `json:"status"`
	Operator string    `json:"operator"`
	CreateAt time.Time `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"autoUpdateTime"`
}
