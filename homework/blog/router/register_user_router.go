package router

import (
	"blog/request"
	"blog/response"
	"blog/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRouter(engine *gin.Engine, db *gorm.DB) {

	userService := service.NewUserService(db)

	userGroup := engine.Group("/user")
	{
		userGroup.GET("/info", handleUserInfo(userService))

		userGroup.POST("/register", handleRegister(userService))

		userGroup.POST("/update", handleUpdate(userService))

		userGroup.DELETE("/delete", handleDelete(userService))
	}

}

func handleUserInfo(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "id 必须是正整数",
			})
			return
		}

		userInfo, err := userService.UserInfoById(&id)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Fail(http.StatusBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, response.Success(request.UserResponse{
			UserName: userInfo.UserName,
			Email:    userInfo.Email,
		}))
	}
}

func handleRegister(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var createRequest request.UserCreateRequest
		err := c.ShouldBindJSON(&createRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Fail(http.StatusBadRequest, err.Error()))
			return
		}

		// 这里还需要做一些校验，比如用户名是否重复，是否有特殊字符等，这里是不是重点
		err = userService.Register(createRequest)

		if err != nil {
			c.JSON(http.StatusBadRequest, response.Fail(http.StatusBadRequest, err.Error()))
			return
		}
		c.JSON(http.StatusOK, response.Success(""))
	}
}

func handleUpdate(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateRequest request.UserUpdateRequest
		err := c.ShouldBindJSON(&updateRequest)
		if err != nil {
			c.JSON(http.StatusOK, response.Fail(http.StatusBadRequest, err.Error()))
			return
		}

		err = userService.Update(updateRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Fail(http.StatusBadRequest, err.Error()))
		}
		c.JSON(http.StatusOK, response.Success(""))
		return
	}
}

func handleDelete(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var deletedRequest request.UserDeletedRequest
		err := c.ShouldBindJSON(&deletedRequest)
		if err != nil {
			c.JSON(http.StatusOK, response.Fail(http.StatusBadRequest, err.Error()))
			return
		}

		err = userService.Delete(deletedRequest)
		if err != nil {
			c.JSON(http.StatusOK, response.Fail(http.StatusBadRequest, err.Error()))
			return
		}
		c.JSON(http.StatusOK, response.Success(""))
		return
	}
}
