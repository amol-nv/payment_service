package payment

import "encoding/json"

type PostPaymentRequest struct {
	CustomerID string  `json:"customerId"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	Reference  string  `json:"reference,omitempty"`
}

type PostPaymentResponse struct {
	PaymentID string          `json:"paymentId"`
	Status    string          `json:"status"`
	Raw       json.RawMessage `json:"-"`
}
