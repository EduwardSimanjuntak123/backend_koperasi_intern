package requests

import "backend_koperasi/internal/models"

type OrderItemRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	Qty       int     `json:"qty" binding:"required,min=1"`
	Notes     *string `json:"notes"`
}

type CreateOrderRequest struct {
	Items          []OrderItemRequest    `json:"items"`
	CourierID      *string               `json:"courier_id"`
	PaymentMethod  *models.PaymentMethod `json:"payment_method"`
	ShippingCost   float64               `json:"shipping_cost"`
	DiscountAmount float64               `json:"discount_amount"`
	Notes          *string               `json:"notes"`
}

type CheckoutCartRequest struct {
	CourierID      *string               `json:"courier_id"`
	PaymentMethod  *models.PaymentMethod `json:"payment_method"`
	ShippingCost   float64               `json:"shipping_cost"`
	DiscountAmount float64               `json:"discount_amount"`
	Notes          *string               `json:"notes"`
}

type UpdateOrderStatusRequest struct {
	Status models.OrderStatus `json:"status" binding:"required"`
	Notes  *string            `json:"notes"`
}

type AssignCourierRequest struct {
	CourierID string `json:"courier_id" binding:"required"`
}

type CancelOrderRequest struct {
	Reason *string `json:"reason"`
}
