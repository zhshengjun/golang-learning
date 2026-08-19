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

func RegisterUserRouter(engine *gin.Engine, db *gorm.DB) {

	userService := service.NewUserService(db)
	engine.POST("/user/register", handleUserRegister(userService))

	loginGroup := engine.Group("/user")
	loginGroup.Use(middleware.RequireLogin())
	{
		loginGroup.GET("/info", handleUserInfo(userService))

		loginGroup.POST("/updated", handleUserUpdated(userService))

		loginGroup.DELETE("/deleted", handleUserDeleted(userService))
	}

}

func handleUserInfo(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, response.Fail(http.StatusBadRequest, "id must be a positive integer"))
			return
		}
		userInfo, err := userService.UserInfoById(&id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response.Success(request.UserResponse{
			UserName: userInfo.UserName,
			Email:    userInfo.Email,
		}))
	}
}

func handleUserRegister(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var createRequest request.UserCreateRequest
		err := c.ShouldBindJSON(&createRequest)
		if err != nil {
			_ = c.Error(fmt.Errorf("%w: param analysis error", blogerrors.ErrBadRequest))
			return
		}

		// 这里还需要做一些校验，比如用户名是否重复，是否有特殊字符等，这里是不是重点
		err = userService.Register(&createRequest)

		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, response.Success(""))
	}
}

func handleUserUpdated(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateRequest request.UserUpdateRequest
		err := c.ShouldBindJSON(&updateRequest)
		if err != nil {
			_ = c.Error(fmt.Errorf("%w: param analysis error", blogerrors.ErrBadRequest))
			return
		}

		name, _ := middleware.CurrentUserName(c)

		updateRequest.Operator = name

		err = userService.Updated(&updateRequest)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, response.Success(""))
		return
	}
}

func handleUserDeleted(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var deletedRequest request.UserDeletedRequest
		err := c.ShouldBindJSON(&deletedRequest)
		if err != nil {
			_ = c.Error(fmt.Errorf("%w: param analysis error", blogerrors.ErrBadRequest))
			return
		}

		name, _ := middleware.CurrentUserName(c)
		deletedRequest.Operator = name

		err = userService.Deleted(&deletedRequest)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, response.Success(""))
		return
	}
}
