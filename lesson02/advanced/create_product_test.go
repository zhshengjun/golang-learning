package advanced

import (
	"database/common"
	"database/models"
	"testing"

	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {

		products := []models.Product{
			{Name: "这是产品01", Category: "类别1", Price: 100},
			{Name: "这是产品02", Category: "类别2", Price: 200},
			{Name: "这是产品03", Category: "类别2", Price: 250},
		}
		db.Create(&products)
	})

}
