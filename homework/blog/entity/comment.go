package entity

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	Id        int64     `json:"id" gorm:"primary_key"`
	ArticleId int64     `json:"articleId" gorm:"index;column:article_id"`
	AnswerId  int64     `json:"answerId" gorm:"index;column:answer_id"` // 回复的ID
	Comment   string    `json:"comment"`
	Status    bool      `json:"status" gorm:"column:status"`
	Creator   string    `json:"creator"`
	CreateAt  time.Time `json:"createAt" gorm:"column:create_at;autoCreateTime"`
	UpdateAt  time.Time `json:"updateAt" gorm:"column:update_at;autoUpdateTime"`
}

func (c *Comment) AfterCreate(db *gorm.DB) error {
	return db.Model(&Article{}).
		Where("id = ?", c.ArticleId).
		UpdateColumns(map[string]any{
			"comment_count":  gorm.Expr("comment_count + ?", 1),
			"comment_status": true,
		}).Error
}

func (c *Comment) AfterUpdate(tx *gorm.DB) error {

	var count int64
	if err := tx.Model(&Comment{}).
		Where("article_id = ? AND status = true", c.ArticleId).
		Count(&count).Error; err != nil {
		return err
	}
	status := false
	if count > 0 {
		status = true
	}

	return tx.Model(&Article{}).
		Where("id = ?", c.ArticleId).
		UpdateColumns(map[string]any{
			"comment_count":  count,
			"comment_status": status,
		}).
		Error
}
