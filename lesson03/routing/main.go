package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Register struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Age      int    `json:"age" binding:"required,gt=0,lt=100"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func main() {

	engine := gin.Default()

	//=====================分组=======================

	v1 := engine.Group("/user/v1")
	{
		v1.POST("/register", register)
		v1.POST("/login", func(c *gin.Context) {})
		v1.POST("/signup", func(c *gin.Context) {})
		v1.POST("/forgot", func(c *gin.Context) {})
		v1.Handle("GET", "/profile", func(c *gin.Context) {})
	}

	v2 := engine.Group("/user/v2")
	{
		v2.POST("/login", login)
	}

	err := engine.Run()
	if err != nil {
		return
	}
}

func register(c *gin.Context) {
	var registerObj Register
	err := c.ShouldBindJSON(&registerObj)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, success(&registerObj))
}

func success(data any) gin.H {
	return gin.H{
		"code":    200,
		"message": "ok",
		"data":    data,
	}
}

func login(context *gin.Context) {
	context.JSON(
		200, gin.H{
			"message": "登录成功",
		})
}
