package model

import "errors"

type PaymentPostRequest struct {
	CustomerID string  `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	Reference  string  `json:"reference,omitempty"`
}

type PaymentPostResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
}

func (r *PaymentPostRequest) Validate() error {
	if r.CustomerID == "" {
		return errors.New("customer_id is required")
	}
	if r.Amount <= 0 {
		return errors.New("amount must be > 0")
	}
	if r.Currency == "" {
		return errors.New("currency is required")
	}
	return nil
}
