package operator

import (
	"database/utils"
	"testing"
)

func TestUpdates(t *testing.T) {

	db := utils.InitDB()

	query := utils.User{ID: 17, Age: 33, Status: true}

	db.Updates(&query)
}
