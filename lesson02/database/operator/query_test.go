package operator

import (
	"database/utils"
	"testing"
)

func TestQuery(t *testing.T) {

	db := utils.InitDB()

	var query utils.User
	db.Where(utils.User{Age: 18}).First(&query)
	utils.FormatePrint(query)
}
