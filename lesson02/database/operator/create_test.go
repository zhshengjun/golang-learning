package operator

import (
	"database/utils"
	"log"
	"testing"
)

func TestCreate(t *testing.T) {

	db := utils.InitDB()

	user := utils.User{
		Name: "张三",
		Age:  18,
	}
	err := db.Create(&user).Error
	if err != nil {
		log.Fatal(err)
	}

	users := []utils.User{
		{Name: "李四", Age: 18},
		{Name: "王五", Age: 19},
		{Name: "周六", Age: 19},
	}
	err = db.Create(&users).Error
	if err != nil {
		log.Fatal(err)
	}

	user = utils.User{Name: "溜", Age: 10}
	err = db.Save(&user).Error
	if err != nil {
		log.Fatal(err)
	}
}
