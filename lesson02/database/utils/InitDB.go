package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func InitDB() (db *gorm.DB) {

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("读取 .env 失败:", err)
	}

	db, err := gorm.Open(
		mysql.Open(os.Getenv("DSN_MYSQL")),
		&gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
