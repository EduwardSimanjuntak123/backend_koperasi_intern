package models

type PaymentMethod struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Category string `gorm:"not null" json:"category"`
	Status   bool   `gorm:"not null" json:"status"`
}
