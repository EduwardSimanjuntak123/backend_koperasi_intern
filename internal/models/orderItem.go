package models

import "time"

type OrderItem struct {
	ID string `gorm:"primaryKey" json:"id"`

	OrderID string `gorm:"not null;index" json:"order_id"`
	Order   Order  `gorm:"foreignKey:OrderID" json:"order,omitempty"`

	ProductID string  `gorm:"not null;index" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`

	Qty int `gorm:"not null" json:"qty"`

	Price float64 `gorm:"not null" json:"price"`

	Discount float64 `gorm:"default:0" json:"discount"`

	Subtotal float64 `gorm:"not null" json:"subtotal"`

	Notes string `gorm:"type:text" json:"notes"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
