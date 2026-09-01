package models

import "time"

type OrderStatusHistory struct {
	ID string `gorm:"primaryKey;size:50" json:"id"`

	OrderID string `gorm:"size:50;not null;index" json:"order_id"`
	Order   *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`

	Status string `gorm:"size:50;not null" json:"status"`

	Description string `gorm:"type:text" json:"description"`

	CreatedBy string `gorm:"size:100" json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
}

