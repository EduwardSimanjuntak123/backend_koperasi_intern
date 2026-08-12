package models

type PaymentMethod struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Category string `gorm:"not null" json:"category"`
	Status   bool   `gorm:"not null" json:"status"`
}
