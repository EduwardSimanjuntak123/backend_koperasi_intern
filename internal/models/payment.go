package models

import "time"

type PaymentMethod string

const (
	PaymentCash     PaymentMethod = "CASH"
	PaymentQRIS     PaymentMethod = "QRIS"
	PaymentTransfer PaymentMethod = "TRANSFER"
	PaymentEWallet  PaymentMethod = "EWALLET"
)

type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "PENDING"
	PaymentPaid     PaymentStatus = "PAID"
	PaymentFailed   PaymentStatus = "FAILED"
	PaymentRefunded PaymentStatus = "REFUNDED"
)

type Payment struct {
	ID string `gorm:"primaryKey;size:50" json:"id"`

	OrderID      string        `gorm:"size:50;uniqueIndex;not null" json:"order_id"`
	Order        *Order        `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Method       PaymentMethod `gorm:"size:30;not null" json:"method"`
	Status       PaymentStatus `gorm:"size:30;default:'PENDING'" json:"status"`
	Amount       float64       `gorm:"not null" json:"amount"`
	PaidAt       *time.Time    `json:"paid_at"`
	ReferenceNo  *string       `gorm:"size:100" json:"reference_no"`
	PaymentProof *string       `gorm:"type:text" json:"payment_proof"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

