package models

import (
	"time"
)

type StoreMember struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	JoinedAt  time.Time `json:"joined_at"`
	Status    string    `gorm:"not null" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	StoreID   string    `gorm:"not null" json:"store_id"`
	Store     Store     `gorm:"foreignKey:StoreID" json:"store"`
	UserID    string    `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	RoleID    string    `gorm:"not null" json:"role_id"`
	Role      Roles     `gorm:"foreignKey:RoleID" json:"role"`
}

func (StoreMember) TableName() string {
	return "store_members"
}
