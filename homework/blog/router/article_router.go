package router

import (
	blogerrors "blog/errors"
	"blog/middleware"
	"blog/request"
	"blog/response"
	"blog/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterArticleRouter(engine *gin.Engine, db *gorm.DB) {

	articleService := service.NewArticleService(db)

	articleGroup := engine.Group("/article")
	articleGroup.Use(middleware.RequireLogin())
	{
		articleGroup.POST("/created", handleArticleCreated(articleService))

		articleGroup.PUT("/updated", handleArticleUpdated(articleService))

		articleGroup.PUT("/published", handleArticlePublished(articleService))

		articleGroup.DELETE("/deleted", handleArticleDeleted(articleService))
	}

}

func handleArticleCreated(articleService *service.ArticleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var createRequest request.ArticleCreateRequest
		err := c.ShouldBindJSON(&createRequest)
		if err != nil {
			_ = c.Error(fmt.Errorf("%w: param analysis error", blogerrors.ErrBadRequest))
			return
		}

		name, _ := middleware.CurrentUserName(c)

		createRequest.Author = name

		// 这里还需要做一些校验，比如用户名是否重复，是否有特殊字符等，这里是不是重点
		err = articleService.Created(&createRequest)

		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, response.Success(""))
	}
}

func handleArticleUpdated(articleService *service.ArticleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateRequest request.ArticleUpdateRequest
		err := c.ShouldBindJSON(&updateRequest)
		if err != nil {
			_ = c.Error(fmt.Errorf("%w: param analysis error", blogerrors.ErrBadRequest))
			return
		}

		name, _ := middleware.CurrentUserName(c)

		updateRequest.Operator = name

		err = articleService.Updated(&updateRequest)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, response.Success(""))
		return
	}
}

func handleArticlePublished(articleService *service.ArticleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var publishedRequest request.ArticlePublishedRequest
		err := c.ShouldBindJSON(&publishedRequest)
		if err != nil {
			_ = c.Error(fmt.Errorf("%w: param analysis error", blogerrors.ErrBadRequest))
			return
		}

		name, _ := middleware.CurrentUserName(c)

		publishedRequest.Operator = name

		err = articleService.Published(&publishedRequest)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, response.Success(""))
		return
	}
}

func handleArticleDeleted(articleService *service.ArticleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var deletedRequest request.ArticleDeleteRequest
		err := c.ShouldBindJSON(&deletedRequest)
		if err != nil {
			_ = c.Error(fmt.Errorf("%w: param analysis error", blogerrors.ErrBadRequest))
			return
		}

		name, _ := middleware.CurrentUserName(c)
		deletedRequest.Operator = name

		err = articleService.Deleted(&deletedRequest)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response.Success(""))
		return
	}
}
