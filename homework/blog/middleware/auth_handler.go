package middleware

import (
	"net/http"

	"blog/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

const loginUserKey = "loginUser"

func RequireLogin() gin.HandlerFunc {

	return func(c *gin.Context) {
		secret := []byte(viper.GetString("jwt.secret"))
		tokenString, err := c.Cookie("blog_token")
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.Fail(http.StatusUnauthorized, "please login first"),
			)
			return
		}

		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (any, error) {
				return secret, nil
			},
			jwt.WithValidMethods([]string{
				jwt.SigningMethodHS256.Alg(),
			}),
			jwt.WithIssuer("blog"),
			jwt.WithExpirationRequired(),
		)
		if err != nil || !token.Valid || claims.Subject == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.Fail(http.StatusUnauthorized, "your login session has expired"),
			)
			return
		}

		// 当前项目的Subject保存的是用户名
		c.Set(loginUserKey, claims.Subject)

		c.Next()
	}
}

func CurrentUserName(c *gin.Context) (string, bool) {
	value, exists := c.Get(loginUserKey)
	if !exists {
		return "", false
	}
	userName, ok := value.(string)
	return userName, ok
}
