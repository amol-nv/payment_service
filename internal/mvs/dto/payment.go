package dto

// PostPaymentRequest is the payload sent to the payment provider.
// Adjust fields to match the payment API contract.
type PostPaymentRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	OrderID  string `json:"order_id"`
	// Optional fields
	CustomerID string `json:"customer_id,omitempty"`
}

// PostPaymentResponse is the expected response from the payment provider.
type PostPaymentResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	// Optional fields
	RedirectURL string `json:"redirect_url,omitempty"`
}

// ErrorResponse represents a typical error shape.
// If the upstream provider uses a different shape, update accordingly.
type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}
