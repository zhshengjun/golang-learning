package entity

import "time"

type Post struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"userId" gorm:"column:user_id;index"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"createAt"`
	UpdateAt  time.Time `json:"updateAt"`
}
