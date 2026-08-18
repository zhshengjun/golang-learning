package main

import "github.com/gin-gonic/gin"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	AuthCode int    `json:"authCode"`
}

func main() {

	engine := gin.Default()

	engine.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World!"})
	})

	engine.GET("/path/:id/:postId", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok",
			"id":     c.Param("id"),
			"postId": c.Param("postId")})
	})

	engine.GET("/path/file/*path", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok",
			"path": c.Param("path")})
	})

	engine.GET("/param", func(c *gin.Context) {
		keyword := c.Query("keyword")
		page := c.DefaultQuery("page", "1")
		pageSize := c.DefaultQuery("pageSize", "20")
		c.JSON(200, gin.H{
			"keyword":  keyword,
			"page":     page,
			"pageSize": pageSize,
		})
	})

	engine.POST("/register/param", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		authCode := c.Query("authCode")

		c.JSON(200, gin.H{
			"username": username,
			"password": password,
			"authCode": authCode,
		})

	})

	engine.POST("/register/body", func(c *gin.Context) {

		var registerRequest RegisterRequest
		err := c.BindJSON(&registerRequest)
		if err != nil {
			return
		}

		c.JSON(200, gin.H{
			"username": registerRequest.Username,
			"password": registerRequest.Password,
			"authCode": registerRequest.AuthCode,
		})

	})

	//=====================分组=======================

	v1 := engine.Group("/user/v1")
	{
		v1.POST("/login", func(c *gin.Context) {})
		v1.POST("/signup", func(c *gin.Context) {})
		v1.POST("/forgot", func(c *gin.Context) {})
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

func login(context *gin.Context) {
	context.JSON(
		200, gin.H{
			"message": "登录成功",
		})
}
