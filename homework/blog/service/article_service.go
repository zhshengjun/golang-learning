package service

import (
	"blog/entity"
	"blog/enums"
	apperrors "blog/errors"
	"blog/request"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ArticleService struct {
	db *gorm.DB
}

func NewArticleService(db *gorm.DB) *ArticleService {
	return &ArticleService{db: db}
}

func (s *ArticleService) Created(createRequest *request.ArticleCreateRequest) error {

	var article entity.Article

	article.Content = createRequest.Content
	article.Title = createRequest.Title
	article.Author = createRequest.Author
	article.Status = enums.ArticleStatusDraft
	article.CommentStatus = false
	article.Operator = createRequest.Author

	result := s.db.Model(&article).Create(&article)

	if result.Error != nil {
		return fmt.Errorf("%w: create article error", result.Error)
	}

	return nil

}

func (s *ArticleService) Updated(updateRequest *request.ArticleUpdateRequest) error {

	article, err := checkArticle(s, updateRequest.Id, &updateRequest.Operator)
	if err != nil {
		return err
	}

	article.Id = updateRequest.Id
	article.Title = updateRequest.Title
	article.Content = updateRequest.Content

	result := s.db.Model(&entity.Article{}).
		Where("id = ?", updateRequest.Id).
		UpdateColumns(map[string]any{
			"title":   updateRequest.Title,
			"content": updateRequest.Content,
		})

	if result.Error != nil {
		return fmt.Errorf("%w: update article error", result.Error)
	}

	return nil
}

func (s *ArticleService) Published(publishedRequest *request.ArticlePublishedRequest) error {

	article, err := checkArticle(s, publishedRequest.Id, &publishedRequest.Operator)
	if err != nil {
		return err
	}

	article.Id = publishedRequest.Id

	result := s.db.Model(&entity.Article{}).
		Where("id = ?", publishedRequest.Id).
		Update("status", enums.ArticleStatusPublished)

	if result.Error != nil {
		return fmt.Errorf("%w: update article error", result.Error)
	}

	return nil
}

func (s *ArticleService) Deleted(deletedRequest *request.ArticleDeleteRequest) error {
	_, err := checkArticle(s, deletedRequest.Id, &deletedRequest.Operator)
	if err != nil {
		return err
	}

	result := s.db.Model(&entity.Article{}).
		Where("id = ?", deletedRequest.Id).
		Update("status", enums.ArticleStatusDeleted)

	if result.Error != nil {
		return fmt.Errorf("%w: update article error", result.Error)
	}
	return nil
}

func (s *ArticleService) UpdateComment(articleId *int64) error {
	_, err := checkArticle(s, *articleId, nil)
	if err != nil {
		return err
	}

	var commentCount int64
	result := s.db.Model(&entity.Comment{}).Where("article_id = ? and status = 0", articleId).Count(&commentCount)

	if result.Error != nil {
		return fmt.Errorf("%w: update comment error", result.Error)
	}

	if commentCount > 0 {
		result := s.db.Model(&entity.Article{}).
			Where("id = ?", articleId).
			UpdateColumns(map[string]any{
				"comment_count":  commentCount,
				"comment_status": true,
			})
		if result.Error != nil {
			return fmt.Errorf("%w: update comment error", result.Error)
		}
	}
	return nil
}

func checkArticle(s *ArticleService, id int64, operator *string) (entity.Article, error) {
	var article entity.Article

	result := s.db.Model(&entity.Article{}).
		Where("id = ?", id).First(&article)
	if &article == nil || article.Id == 0 || errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return entity.Article{}, fmt.Errorf("%w: article not found", apperrors.ErrNotFound)
	}

	if result.Error != nil {
		return entity.Article{}, fmt.Errorf("%w: article error", result.Error)
	}

	if operator != nil {
		if article.Author != *operator {
			return entity.Article{}, fmt.Errorf("%w: update article permission denied", apperrors.ErrForbidden)
		}
	}
	return article, nil
}
