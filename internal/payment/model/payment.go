package model

type PaymentPostRequest struct {
	// Minimal, generic fields; adjust as needed for the external contract.
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	OrderID  string `json:"order_id"`
	// Optional metadata
	CustomerID string `json:"customer_id,omitempty"`
}

type PaymentPostResponse struct {
	Status  string `json:"status"`
	TxnID   string `json:"txn_id"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
