package middleware

import (
	blogerrors "blog/errors"
	"blog/response"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err
		fmt.Println(err)
		status, message := statusOf(err)

		if status >= 500 {
			slog.Error("request failed", "error", err)
		}

		c.AbortWithStatusJSON(
			status,
			response.Fail(status, message),
		)
	}
}

func statusOf(err error) (int, string) {
	switch {
	case errors.Is(err, blogerrors.ErrBadRequest):
		return http.StatusBadRequest, "请求参数错误"
	case errors.Is(err, blogerrors.ErrUnauthorized):
		return http.StatusUnauthorized, "用户认证失败"
	case errors.Is(err, blogerrors.ErrForbidden):
		return http.StatusForbidden, "没有操作权限"
	case errors.Is(err, blogerrors.ErrNotFound):
		return http.StatusNotFound, "文章、评论或用户不存在"
	case errors.Is(err, blogerrors.ErrConflict):
		return http.StatusConflict, "数据已存在"
	case errors.Is(err, blogerrors.ErrGone):
		return http.StatusGone, "数据状态已变更"
	default:
		return http.StatusInternalServerError, "服务器内部错误"
	}
}
