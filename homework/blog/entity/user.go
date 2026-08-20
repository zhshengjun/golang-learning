package entity

import (
	"blog/enums"
	"time"
)

type User struct {
	Id           int64            `json:"id" gorm:"primary_key"`
	UserName     string           `json:"userName" gorm:"column:user_name"`
	Password     string           `json:"-"`
	Email        string           `json:"email"`
	ArticleNum   int              `json:"articleNum" gorm:"column:article_num"`
	Status       enums.UserStatus `json:"status"`
	Certificated bool             `json:"certificated"`
	Operator     string           `json:"operator"`
	CreateAt     time.Time        `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt     time.Time        `json:"updateAt" gorm:"autoUpdateTime"`
}
