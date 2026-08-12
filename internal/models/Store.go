package models

import "time"

type Store struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	NameStore   string    `gorm:"not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	Logo        string    `gorm:"not null" json:"logo"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
}

func (Shop) TableName() string {
	return "shop"
}