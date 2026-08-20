package service

import (
	"blog/entity"
	"blog/enums"
	blogerrors "blog/errors"
	"blog/request"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) Created(createRequest *request.CommentCreateRequest) error {
	var article entity.Article
	result := s.db.
		Where("id = ? AND status != ?", createRequest.ArticleId, enums.ArticleStatusDeleted).
		First(&article)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: article not found", blogerrors.ErrNotFound)
		}
		return fmt.Errorf("%w: query article error", result.Error)
	}
	if article.Id == 0 {
		return fmt.Errorf("%w: article not found", blogerrors.ErrNotFound)
	}

	comment := entity.Comment{
		ArticleId: createRequest.ArticleId,
		AnswerId:  createRequest.AnswerId,
		Comment:   createRequest.Comment,
		Status:    true,
		Creator:   createRequest.Creator,
	}
	if result := s.db.Create(&comment); result.Error != nil {
		return fmt.Errorf("%w: create comment error", result.Error)
	}
	return nil
}

func (s *CommentService) Deleted(deleteRequest *request.CommentDeleteRequest) error {
	comment, err := s.findInfo(deleteRequest.Id)
	if err != nil {
		return err
	}
	if comment.Creator != deleteRequest.Operator {
		return fmt.Errorf("%w: delete comment permission denied", blogerrors.ErrForbidden)
	}

	if result := s.db.Model(&comment).
		Or("answer_id = ?", deleteRequest.Id).
		Update("status", false); result.Error != nil {
		return fmt.Errorf("%w: delete comment error", result.Error)
	}

	return nil
}

func (s *CommentService) CommentList(articleId int) ([]entity.Comment, error) {
	comments := make([]entity.Comment, 0)
	if err := s.db.Model(&entity.Comment{}).
		Where("article_id = ? AND status = ?", articleId, true).
		Order("create_at ASC").Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("%w: query comments error", err)
	}
	return comments, nil
}

func (s *CommentService) findInfo(id int) (entity.Comment, error) {
	var comment entity.Comment
	result := s.db.Model(&entity.Comment{}).Where("id = ? AND status = ?", id, true).First(&comment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || comment.Id == 0 {
		return comment, fmt.Errorf("%w: comment not found", blogerrors.ErrNotFound)
	}
	if result.Error != nil {
		return comment, fmt.Errorf("%w: query comment error", result.Error)
	}
	return comment, nil
}
