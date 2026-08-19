package entity

import "time"

type Article struct {
	ID       int64     `json:"id" gorm:"primary_key"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Author   string    `json:"author" gorm:"column:author;index"`
	Operator string    `json:"operator"`
	CreateAt time.Time `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"autoUpdateTime"`
}
