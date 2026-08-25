package common

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func ExecSql(f func(*gorm.DB)) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("读取 .env 失败:", err)
	}

	db, err := gorm.Open(
		mysql.Open(os.Getenv("DSN_MYSQL")),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   gormLogger.Default.LogMode(gormLogger.Info),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	// 执行具体的甘肃
	f(db)

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
