package models

import "time"

type Brand struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Slug string `gorm:"size:120;uniqueIndex;not null" json:"slug"`

	Logo *string `json:"logo"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	Products []Product `gorm:"foreignKey:BrandID" json:"products,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
