package operator

import (
	"database/utils"
	"testing"
)

func TestQueryPluck(t *testing.T) {

	db := utils.InitDB()
	var names []string
	db.Model(&utils.User{}).Where(utils.User{Age: 20}).
		Pluck("name", &names)

	utils.FormatePrint(names)
}
