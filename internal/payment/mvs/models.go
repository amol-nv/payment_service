package mvs

// PaymentPostRequest is the MVS Model for the outbound POST call.
// It represents the data the payment service will send to an external payment endpoint.
type PaymentPostRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	OrderID  string `json:"order_id"`
}

// PaymentPostResponse is the MVS Model for the expected response from the external payment endpoint.
type PaymentPostResponse struct {
	Status    string `json:"status"`
	PaymentID string `json:"payment_id"`
	Message   string `json:"message,omitempty"`
}
