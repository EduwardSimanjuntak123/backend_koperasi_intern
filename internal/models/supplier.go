package models

import "time"

type Supplier struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:50;uniqueIndex" json:"code"`
	Name string `gorm:"size:150;not null" json:"name"`

	ContactName string `gorm:"size:100" json:"contact_name"`
	Phone       string `gorm:"size:30" json:"phone"`
	Email       string `gorm:"size:100" json:"email"`
	Address     string `gorm:"type:text" json:"address"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	Products []Product `gorm:"foreignKey:SupplierID" json:"products,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
