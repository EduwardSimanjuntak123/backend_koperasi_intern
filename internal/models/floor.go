package models

import "time"

type Floor struct {
	ID   string `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:100;not null" json:"name"`

	DeliveryFee float64 `gorm:"default:0" json:"delivery_fee"`

	BuildingID string   `gorm:"not null" json:"building_id"`
	Building   Building `gorm:"foreignKey:BuildingID" json:"building"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
