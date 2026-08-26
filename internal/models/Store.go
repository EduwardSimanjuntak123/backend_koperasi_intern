package models

import "time"

type Store struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	NameStore   string    `gorm:"not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	Logo        string    `gorm:"not null" json:"logo"`
	UserID      string    `gorm:"not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Store) TableName() string {
	return "stores"
}
