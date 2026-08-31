package models

import "time"

type Courier struct {
	ID string `gorm:"primaryKey" json:"id"`

	Name        string `gorm:"size:100;not null" json:"name"`
	PhoneNumber string `gorm:"size:20" json:"phone_number"`
	IsActive    *bool  `gorm:"default:true" json:"is_active"`

	Orders []Order `gorm:"foreignKey:CourierID" json:"orders,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
