package models

import "time"

type StatusCount struct {
	Status bool  `json:"status" gorm:"column:status"`
	Count  int64 `json:"count" gorm:"column:count"`
}
type User struct {
	ID        uint   `gorm:"primaryKey"  json:"id"`
	Name      string `gorm:"size:255"  json:"name"`
	Age       int    `gorm:"default:0"  json:"age"`
	Status    bool   `gorm:"default:false"  json:"status"`
	Profile   Profile
	Orders    []Order
	Roles     []Role    `gorm:"many2many:user_roles;"`
	CreatedAt time.Time `gorm:"autoCreateTime"  json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"  json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type Profile struct {
	ID        uint      `gorm:"primaryKey"  json:"id"`
	UserID    uint      `gorm:"uniqueIndex"  json:"UserId"`
	NikeName  string    `gorm:"size:255" json:"nikeName"`
	CreatedAt time.Time `gorm:"autoCreateTime"  json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"  json:"updated_at"`
}

type Order struct {
	ID         uint `gorm:"primaryKey"  json:"id"`
	OrderID    uint `gorm:"uniqueIndex"  json:"orderId"`
	UserID     uint `json:"userId"`
	OrderItems []OrderItem
	CreatedAt  time.Time `gorm:"autoCreateTime"  json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"  json:"updated_at"`
}

type OrderItem struct {
	ID          uint      `gorm:"primaryKey"  json:"id"`
	OrderID     uint      `json:"orderId"`
	ProductID   uint      `json:"productId"`
	Product     Product   `json:"product"`
	Quantity    uint      `json:"quantity"`
	TotalAmount uint      `json:"TotalAmount"`
	CreatedAt   time.Time `gorm:"autoCreateTime"  json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"  json:"updated_at"`
}

type Product struct {
	ID        uint      `gorm:"primaryKey"  json:"id"`
	Name      string    `gorm:"size:255"  json:"name"`
	Category  string    `gorm:"size:255"  json:"type"`
	Price     uint      `json:"price"`
	CreatedAt time.Time `gorm:"autoCreateTime"  json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"  json:"updated_at"`
}

type Role struct {
	ID        uint      `gorm:"primaryKey"  json:"id"`
	RoleName  string    `gorm:"size:255"  json:"roleName"`
	Users     []User    `gorm:"many2many:user_roles;"`
	CreatedAt time.Time `gorm:"autoCreateTime"  json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"  json:"updated_at"`
}

type UserRole struct {
	ID        uint      `gorm:"primaryKey"  json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:uk_user_role"`
	RoleID    uint      `gorm:"not null;uniqueIndex:uk_user_role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
