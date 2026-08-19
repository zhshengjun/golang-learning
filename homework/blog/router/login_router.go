package router

import (
	blogerrors "blog/errors"
	"blog/middleware"
	"blog/request"
	"blog/response"
	"blog/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func RegisterLoginRouter(engine *gin.Engine, db *gorm.DB) {

	userService := service.NewUserService(db)
	loginService := service.NewLoginService(db, userService)

	engine.POST("/login", handLogin(loginService))

	loginGroup := engine.Group("")
	loginGroup.Use(middleware.RequireLogin())
	{
		loginGroup.POST("/logout", handLogout())
	}

}

func handLogin(loginService *service.LoginService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest request.LoginRequest
		err := c.ShouldBindJSON(&loginRequest)
		if err != nil {
			_ = c.Error(fmt.Errorf("%w: param analysis error", blogerrors.ErrBadRequest))
			return
		}
		login, err := loginService.VerifyPassword(&loginRequest)
		if err != nil {
			_ = c.Error(err)
			return
		}
		if !login {
			_ = c.Error(fmt.Errorf("%w: login error", blogerrors.ErrBadRequest))
			return
		}

		// 密码验证整ok，创建jwt，保存到cookie
		secret := []byte(viper.GetString("jwt.secret"))

		now := time.Now()
		claims := jwt.RegisteredClaims{
			Subject:   loginRequest.UserName,
			Issuer:    "blog",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour)),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		signedString, err := token.SignedString(secret)
		fmt.Printf("登录的签名: %s\n", signedString)
		setLoginCookie(c, signedString)

		c.JSON(http.StatusOK, response.Success("登录成功"))
	}
}

func handLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		deleteLoginCookie(c)
		c.JSON(http.StatusOK, response.Success("登出成功"))
	}
}

func setLoginCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		"blog_token", // Cookie名称
		token,        // JWT
		2*60*60,      // 有效期：2小时，单位秒
		"/",          // 整个网站都可以携带
		"",           // 当前域名
		false,        // 本地HTTP用false；生产HTTPS用true
		true,         // HttpOnly，禁止JS读取
	)
}

func deleteLoginCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		"blog_token",
		"",
		-1,  // 立即删除
		"/", // 必须和创建时一致
		"",
		false, // 生产HTTPS改成true
		true,
	)
}
