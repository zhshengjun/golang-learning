package middleware

import (
	apperrors "blog/errors"
	"blog/response"
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err
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
	case errors.Is(err, apperrors.ErrBadRequest):
		return 400, "请求参数错误"
	case errors.Is(err, apperrors.ErrUnauthorized):
		return 401, "用户认证失败"
	case errors.Is(err, apperrors.ErrForbidden):
		return 403, "没有操作权限"
	case errors.Is(err, apperrors.ErrNotFound):
		return 404, "文章、评论或用户不存在"
	case errors.Is(err, apperrors.ErrConflict):
		return 409, "数据已存在"
	case errors.Is(err, apperrors.ErrGone):
		return 410, "数据状态已变更"
	default:
		return 500, "服务器内部错误"
	}
}
