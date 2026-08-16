package operator

import (
	"database/utils"
	"testing"
)

func TestQueryOffset(t *testing.T) {

	db := utils.InitDB()

	var query_list []utils.User
	db.Where(utils.User{Age: 20}).Offset(5).Limit(2).Find(&query_list)
	utils.FormatePrint(query_list)
}
