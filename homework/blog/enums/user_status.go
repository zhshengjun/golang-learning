package enums

type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"    // 活跃
	UserStatusCancelled UserStatus = "CANCELLED" // 已注销
)
