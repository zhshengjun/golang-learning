package operator

import (
	"database/utils"
	"testing"
)

func TestUpdateOne(t *testing.T) {

	db := utils.InitDB()

	query := utils.User{ID: 17}

	db.Model(&query).Update("age", 44)
}
