package models

import (
	"time"
)

type Roles struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	StoreID   uint      `gorm:"not null" json:"store_id"`
	Store     Store     `gorm:"foreignKey:StoreID" json:"store"`
}
