package models

import "time"

type Banner struct {
	ID string `gorm:"primaryKey" json:"id"`

	Title string `gorm:"size:150" json:"title"`

	Image string `gorm:"not null" json:"image"`

	Link string `json:"link"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	Order int `gorm:"default:0" json:"order"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
