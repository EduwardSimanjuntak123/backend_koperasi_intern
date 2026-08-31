package models

type DiscountProduct struct {
	DiscountID string `gorm:"primaryKey"`
	ProductID  string `gorm:"primaryKey"`
}
