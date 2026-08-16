package operator

import (
	"database/utils"
	"testing"
)

func TestUpdateColumns(t *testing.T) {

	db := utils.InitDB()

	query := utils.User{ID: 19}

	db.Model(&query).UpdateColumns(map[string]any{
		"age":    89,
		"status": true,
	})
}
