package models

import "time"

type Unit struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name   string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Symbol string `gorm:"size:20;uniqueIndex;not null" json:"symbol"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	Products []Product `gorm:"foreignKey:UnitID" json:"products,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
