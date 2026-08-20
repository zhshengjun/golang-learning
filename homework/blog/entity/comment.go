package entity

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	Id        int       `json:"id" gorm:"primary_key"`
	ArticleId int       `json:"articleId" gorm:"index;column:article_id"`
	AnswerId  int       `json:"answerId" gorm:"index;column:answer_id"` // 回复的ID
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
		Where("article_id = ? AND status = ?", c.ArticleId, true).
		Count(&count).Error; err != nil {
		return err
	}

	return tx.Model(&Article{}).
		Where("id = ?", c.ArticleId).
		UpdateColumns(map[string]any{
			"comment_count":  count,
			"comment_status": count > 0,
		}).Error
}
