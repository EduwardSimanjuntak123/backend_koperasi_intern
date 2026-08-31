package models

import "time"

type OrderStatus string

const (
	OrderPending   OrderStatus = "PENDING"
	OrderPaid      OrderStatus = "PAID"
	OrderPacking   OrderStatus = "PACKING"
	OrderShipping  OrderStatus = "SHIPPING"
	OrderDelivered OrderStatus = "DELIVERED"
	OrderCompleted OrderStatus = "COMPLETED"
	OrderCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID string `gorm:"primaryKey" json:"id"`
	// Customer
	UserID string `gorm:"not null" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"user"`

	CourierID *string  `json:"courier_id"`
	Courier   *Courier `gorm:"foreignKey:CourierID" json:"courier"`

	// Payment
	Payment *Payment `gorm:"foreignKey:OrderID" json:"payment,omitempty"`

	// Order Information
	InvoiceNumber string      `gorm:"size:50;uniqueIndex" json:"invoice_number"`
	Status        OrderStatus `gorm:"size:30;default:'PENDING'" json:"status"`
	Notes         *string     `json:"notes"`
	// Price Summary
	Subtotal       float64 `gorm:"not null;default:0" json:"subtotal"`
	DiscountAmount float64 `gorm:"default:0" json:"discount_amount"`
	ShippingCost   float64 `gorm:"default:0" json:"shipping_cost"`

	GrandTotal float64 `gorm:"not null;default:0" json:"grand_total"`

	// Time
	PaidAt      *time.Time `json:"paid_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
	CancelledAt *time.Time `json:"cancelled_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Items         []OrderItem          `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	StatusHistory []OrderStatusHistory `gorm:"foreignKey:OrderID" json:"status_history,omitempty"`
}
