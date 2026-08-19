package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(engine *gin.Engine, db *gorm.DB) {
	RegisterLoginRouter(engine, db) //登录接口
	RegisterUserRouter(engine, db)  // 用户接口
	RegisterArticleRouter(engine, db)
}
