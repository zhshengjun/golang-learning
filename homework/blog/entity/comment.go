package entity

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	ArticleID uint      `json:"articleId" gorm:"index;column:article_id"`
	Comment   string    `json:"comment"`
	Creator   string    `json:"creator"`
	CreateAt  time.Time `json:"createAt" gorm:"column:create_at;autoCreateTime"`
	UpdateAt  time.Time `json:"updateAt" gorm:"column:update_at;autoUpdateTime"`
}
