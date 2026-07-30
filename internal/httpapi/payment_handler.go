package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"amol-nv/payment_service/internal/payment"
	"amol-nv/payment_service/internal/service"
)

type PaymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(s service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

func (h *PaymentHandler) PostPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req payment.PostPaymentRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "invalid request", "details": err.Error()})
		return
	}

	if strings.TrimSpace(req.CustomerID) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "customerId is required"})
		return
	}
	if req.Amount <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "amount must be > 0"})
		return
	}
	if strings.TrimSpace(req.Currency) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "currency is required"})
		return
	}

	resp, err := h.service.PostPayment(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, contextCanceledOrDeadlineExceeded(err)) {
			status = http.StatusGatewayTimeout
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "payment failed", "details": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// contextCanceledOrDeadlineExceeded is a small helper to avoid importing context in handler.
func contextCanceledOrDeadlineExceeded(err error) error {
	// Best-effort mapping; service already wraps errors.
	if err == nil {
		return nil
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "context canceled") || strings.Contains(msg, "deadline exceeded") {
		return err
	}
	return nil
}
