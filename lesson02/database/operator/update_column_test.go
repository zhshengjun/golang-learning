package operator

import (
	"database/utils"
	"testing"
)

func TestUpdateColumn(t *testing.T) {

	db := utils.InitDB()

	query := utils.User{ID: 18}

	db.Model(&query).UpdateColumn("age", 66)
}
