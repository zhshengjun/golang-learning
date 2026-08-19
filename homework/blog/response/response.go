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

func Success(data any) gin.H {
	return gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data":    data,
	}
}

func Exception(data any) gin.H {
	return gin.H{
		"code":    http.StatusBadRequest,
		"message": "",
		"data":    data,
	}
}

func Fail(code int, message string) gin.H {
	return gin.H{
		"code":    code,
		"message": message,
	}
}
