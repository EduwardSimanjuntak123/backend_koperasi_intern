package models

import "time"

type DiscountType string

const (
	DiscountPercent DiscountType = "PERCENT"
	DiscountNominal DiscountType = "NOMINAL"
)

type Discount struct {
	ID string `gorm:"primaryKey" json:"id"`

	Name string `gorm:"size:150;not null" json:"name"`
	Slug string `gorm:"size:150;uniqueIndex" json:"slug"`

	Type  DiscountType `gorm:"not null" json:"type"`
	Value float64      `gorm:"not null" json:"value"`

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	Products []Product `gorm:"many2many:discount_products;" json:"products,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
