package main

import (
	"fmt"
	"net/http"
	"time"

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

func loggerWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		timeStart := time.Now()
		fmt.Printf("开始处理接口 %s\n", timeStart)
		c.Next()
		fmt.Printf("接口处理完成，耗时 %s\n", time.Since(timeStart))
	}

}

func main() {

	engine := gin.Default()
	// 全局中间件
	//engine.Use(loggerWare())

	//=====================分组=======================

	v1 := engine.Group("/user/v1")
	// 分组使用中间件
	v1.Use(loggerWare())
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

func success(data any) Response {
	return Response{
		Code:    200,
		Message: "ok",
		Data:    data,
	}
}

func login(context *gin.Context) {
	context.JSON(
		200, gin.H{
			"message": "登录成功",
		})
}
