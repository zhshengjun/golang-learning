package operator

import (
	"database/utils"
	"testing"
)

func TestDelete(t *testing.T) {

	// 直接使用 Delete(&user)，需要 user 中有主键
	db := utils.InitDB()
	user := utils.User{Name: "张三", ID: 15}
	db.Delete(&user)

	// 条件没有主见，必须使用Where，而且在 Delete之前
	user = utils.User{Name: "李四"}
	db.Where(utils.User{Name: "李四"}).Delete(&user)
}
