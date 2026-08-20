package service

import (
	"blog/entity"
	"blog/enums"
	blogerrors "blog/errors"
	"blog/request"
	"blog/response"
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

	var user entity.User
	result := s.db.Model(&entity.User{}).Where("user_name = ? and certificated = ?", createRequest.Author, true).Find(&user)

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: author not certificated or not exist", blogerrors.ErrForbidden)
	}

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: author not certificated", blogerrors.ErrNotFound)
		}
		return fmt.Errorf("%w: user query error", result.Error)
	}

	var article entity.Article

	article.Content = createRequest.Content
	article.Title = createRequest.Title
	article.Author = createRequest.Author
	article.Status = enums.ArticleStatusDraft
	article.CommentStatus = false
	article.Operator = createRequest.Author

	result = s.db.Model(&entity.Article{}).Create(&article)

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

	if article.Status == enums.ArticleStatusDeleted {
		return fmt.Errorf("%w:  article has deleted", blogerrors.ErrGone)
	}

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

	result := s.db.Model(&article).
		Where("id = ?", publishedRequest.Id).
		Update("status", enums.ArticleStatusPublished)

	if result.Error != nil {
		return fmt.Errorf("%w: update article error", result.Error)
	}

	return nil
}

func (s *ArticleService) Deleted(deletedRequest *request.ArticleDeleteRequest) error {
	article, err := checkArticle(s, deletedRequest.Id, &deletedRequest.Operator)
	if err != nil {
		return err
	}

	result := s.db.Model(&article).
		Update("status", enums.ArticleStatusDeleted)

	if result.Error != nil {
		return fmt.Errorf("%w: update article error", result.Error)
	}
	return nil
}

func (s *ArticleService) ArticleInfo(articleId *int) (entity.Article, error) {
	article, err := checkArticle(s, *articleId, nil)
	if err != nil {
		return entity.Article{}, err
	}
	return *article, nil
}

func (s *ArticleService) ArticlePageList(articleListRequest *request.ArticleListRequest) (*response.PageResponse[entity.Article], error) {

	var total int64
	if err := s.db.Model(&entity.Article{}).Where("author = ? and status != ?", articleListRequest.Author, enums.ArticleStatusDeleted).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("%w: query articles error", err)
	}
	responseArticle := &response.PageResponse[entity.Article]{
		Total: int(total),
		Pages: Pages(int(total), articleListRequest.PageSize),
	}

	var articles []entity.Article
	if err := s.db.Model(&entity.Article{}).Scopes(Paginate(articleListRequest.CurrentPage, articleListRequest.PageSize)).
		Where("author = ? and status != ?", articleListRequest.Author, enums.ArticleStatusDeleted).
		Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("%w: query articles error", err)
	}
	responseArticle.List = articles

	return responseArticle, nil

}

func checkArticle(s *ArticleService, id int, operator *string) (*entity.Article, error) {
	var article entity.Article

	result := s.db.Model(&entity.Article{}).Where("id = ?", id).First(&article)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: article not found", blogerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: article query error", result.Error)
	}

	if article.Id == 0 {
		return nil, fmt.Errorf("%w: article not found", blogerrors.ErrNotFound)
	}

	if article.Status == enums.ArticleStatusDeleted {
		return nil, fmt.Errorf("%w:  article is deleted", blogerrors.ErrGone)
	}

	if operator != nil {
		if article.Author != *operator {
			return nil, fmt.Errorf("%w: update article permission denied", blogerrors.ErrForbidden)
		}
	}
	return &article, nil
}
