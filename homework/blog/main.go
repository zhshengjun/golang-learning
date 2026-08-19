package main

import (
	"blog/router"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func init() {
	viper.AutomaticEnv()
	viper.SetConfigFile("config.yaml")
	viper.SetEnvPrefix("BLOG")
	err := viper.ReadInConfig()
	if err != nil {
		return
	}
}

func main() {
	engine := gin.Default()
	db := initDB()
	defer closeDB(db)

	// 这里注册路由
	router.Register(engine, db)

	err := engine.Run()
	if err != nil {
		return
	}
}

func initDB() *gorm.DB {
	db, err := gorm.Open(
		mysql.Open(viper.GetString("database.dsn")),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   gormLogger.Default.LogMode(gormLogger.Info),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func closeDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(sqlDB)
}
