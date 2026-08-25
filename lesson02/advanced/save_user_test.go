package advanced

import (
	"database/common"
	"database/models"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestSaveUser(t *testing.T) {

	common.ExecSql(func(db *gorm.DB) {

		user := models.User{
			ID: 30, Name: "张三", Age: 34, Status: true,
			CreatedAt: time.Now(),
			Profile: models.Profile{
				NikeName: "八星古神",
			},
			Orders: []models.Order{
				{OrderID: 1, UserID: 30,
					OrderItems: []models.OrderItem{
						{OrderID: 1, ProductID: 1, Quantity: 1, TotalAmount: 120},
						{OrderID: 1, ProductID: 2, Quantity: 1, TotalAmount: 180},
					},
				},
				{OrderID: 2, UserID: 30,
					OrderItems: []models.OrderItem{
						{OrderID: 2, ProductID: 3, Quantity: 1, TotalAmount: 160},
					},
				},
			},
			//Roles: []models.Role{
			//	{RoleName: "角色01"},
			//},
		}
		db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user)
	})

}
