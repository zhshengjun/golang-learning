package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(engine *gin.Engine, db *gorm.DB) {
	RegisterUserRouter(engine, db)
}
