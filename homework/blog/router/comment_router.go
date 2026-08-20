package router

import (
	blogerrors "blog/errors"
	"blog/middleware"
	"blog/request"
	"blog/response"
	"blog/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCommentRouter(engine *gin.Engine, db *gorm.DB) {
	commentService := service.NewCommentService(db)

	commentGroup := engine.Group("/comment")
	commentGroup.Use(middleware.RequireLogin())
	{
		commentGroup.POST("/created", handleCommentCreated(commentService))
		commentGroup.DELETE("/deleted", handleCommentDeleted(commentService))
		commentGroup.GET("/list", handleCommentList(commentService))
	}
}

func handleCommentCreated(commentService *service.CommentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var createRequest request.CommentCreateRequest
		if err := c.ShouldBindJSON(&createRequest); err != nil {
			_ = c.Error(fmt.Errorf("%w: param analysis error", blogerrors.ErrBadRequest))
			return
		}

		createRequest.Creator, _ = middleware.CurrentUserName(c)
		if err := commentService.Created(&createRequest); err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, response.Success(""))
	}
}

func handleCommentDeleted(commentService *service.CommentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var deleteRequest request.CommentDeleteRequest
		if err := c.ShouldBindJSON(&deleteRequest); err != nil {
			_ = c.Error(fmt.Errorf("%w: param analysis error", blogerrors.ErrBadRequest))
			return
		}

		deleteRequest.Operator, _ = middleware.CurrentUserName(c)
		if err := commentService.Deleted(&deleteRequest); err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, response.Success(""))
	}
}

func handleCommentList(commentService *service.CommentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		articleId, err := strconv.Atoi(c.Query("articleId"))
		if err != nil || articleId <= 0 {
			c.JSON(http.StatusBadRequest, response.Fail(http.StatusBadRequest, "articleId must be a positive integer"))
			return
		}

		comments, err := commentService.CommentList(articleId)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, response.Success(comments))
	}
}
