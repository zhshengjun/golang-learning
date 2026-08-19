package enums

type UserStatus string

const (
	UserStatusActive     UserStatus = "ACTIVE"     // 活跃
	UserStatusFrozen     UserStatus = "FROZEN"     // 冻结
	UserStatusCancelling UserStatus = "CANCELLING" // 注销中
	UserStatusCancelled  UserStatus = "CANCELLED"  // 已注销
)
