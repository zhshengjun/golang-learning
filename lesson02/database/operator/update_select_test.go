package operator

import (
	"database/utils"
	"testing"
)

func TestUpdateSelect(t *testing.T) {

	db := utils.InitDB()

	query := utils.User{ID: 20, Name: "test", Age: 51, Status: true}

	db.Model(&query).Select("age", "status").Updates(&query)
}
