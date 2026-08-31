package models

import "time"

type VoucherType string

const (
	VoucherPercent VoucherType = "PERCENT"
	VoucherNominal VoucherType = "NOMINAL"
)

type Voucher struct {
	ID string `gorm:"primaryKey" json:"id"`

	Code string `gorm:"uniqueIndex;not null" json:"code"`

	Name string `gorm:"not null" json:"name"`

	Type  VoucherType `json:"type"`
	Value float64     `json:"value"`

	MinimumTransaction float64 `json:"minimum_transaction"`

	MaxDiscount float64 `json:"max_discount"`

	Quota int `json:"quota"`

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
