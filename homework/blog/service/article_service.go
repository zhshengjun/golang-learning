package service

import (
	"blog/entity"
	"blog/request"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ArticleService struct {
	DB *gorm.DB
}

func NewArticleService(db *gorm.DB) *ArticleService {
	return &ArticleService{DB: db}
}

func (s *ArticleService) Create(createRequest *request.ArticleCreateRequest) {

	var article entity.Article

	article.Content = createRequest.Content
	article.Title = createRequest.Title
	article.Author = createRequest.Author
	article.Operator = createRequest.Author
	article.CreateAt = time.Now()
	article.UpdateAt = time.Now()

	s.DB.Model(&article).Create(&article)

}

func (s *ArticleService) Update(updateRequest *request.ArticleUpdateRequest) error {

	var article entity.Article

	s.DB.Model(&article).
		Where("id = ?", updateRequest.ID).First(&article)

	if article.ID == 0 || article.Author != updateRequest.Operator {
		return errors.New("permission denied")
	}

	article.ID = updateRequest.ID
	article.Title = updateRequest.Title

	return nil
}

func (s *ArticleService) Delete(updateRequest *request.ArticleDeleteRequest) {

}
