package models

import (
	"time"
)

type StoreMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	joinedAt  time.Time `json:"joined_at"`
	status    string    `gorm:"not null" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	StoreID   uint      `gorm:"not null" json:"store_id"`
	Store     Store     `gorm:"foreignKey:StoreID" json:"store"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}
