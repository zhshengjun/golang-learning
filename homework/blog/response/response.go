package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type PageResponse[T any] struct {
	Total int `json:"total"`
	Pages int `json:"pages"`
	List  []T `json:"list"`
}

func Success(data any) gin.H {
	return gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data":    data,
	}
}

func Fail(code int, message string) gin.H {
	return gin.H{
		"code":    code,
		"message": message,
	}
}
