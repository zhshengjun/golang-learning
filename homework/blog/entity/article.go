package entity

import (
	"blog/enums"
	"time"

	"gorm.io/gorm"
)

type Article struct {
	Id            int                 `json:"id" gorm:"primary_key"`
	Title         string              `json:"title"`
	Content       string              `json:"content"`
	Author        string              `json:"author" gorm:"column:author;index"`
	Status        enums.ArticleStatus `json:"status" gorm:"column:status"`
	CommentStatus bool                `json:"CommentStatus" gorm:"column:comment_status"`
	CommentCount  int                 `json:"commentCount" gorm:"column:comment_count"`
	Operator      string              `json:"operator"`
	CreateAt      time.Time           `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt      time.Time           `json:"updateAt" gorm:"autoUpdateTime"`
}

func (a *Article) AfterUpdate(tx *gorm.DB) error {

	var count int64
	if err := tx.Model(&Article{}).
		Where("author = ? AND status = 'PUBLISHED'", a.Author).
		Count(&count).Error; err != nil {
		return err
	}

	return tx.Model(&User{}).
		Where("user_name = ?", a.Author).
		UpdateColumn("article_num", count).Error
}
