package model

import "time"

type PaymentRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Method   string `json:"method"`
	OrderID  string `json:"order_id"`
	// Optional metadata
	CustomerID string `json:"customer_id,omitempty"`
	// Optional idempotency key
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}

type PaymentResponse struct {
	PaymentID string    `json:"payment_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
