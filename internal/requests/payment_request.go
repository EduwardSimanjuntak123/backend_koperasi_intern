package requests

import "backend_koperasi/internal/models"

type CreatePaymentRequest struct {
	OrderID      string               `json:"order_id" binding:"required"`
	Method       models.PaymentMethod `json:"method" binding:"required"`
	Amount       float64              `json:"amount" binding:"required,gt=0"`
	ReferenceNo  *string              `json:"reference_no"`
	PaymentProof *string              `json:"payment_proof"`
}

type PayOrderRequest struct {
	OrderID      string                `json:"order_id" binding:"required"`
	Method       *models.PaymentMethod `json:"method"`
	ReferenceNo  *string               `json:"reference_no"`
	PaymentProof *string               `json:"payment_proof"`
}

type UpdatePaymentRequest struct {
	Method       *models.PaymentMethod `json:"method"`
	Status       *models.PaymentStatus `json:"status"`
	Amount       *float64              `json:"amount"`
	ReferenceNo  *string               `json:"reference_no"`
	PaymentProof *string               `json:"payment_proof"`
}
