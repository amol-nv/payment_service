package view

import (
	"encoding/json"
	"net/http"

	"amol-nv/payment_service/internal/mvs/dto"
	"amol-nv/payment_service/internal/payment/service"
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

	var req dto.PostPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Message: "invalid json"})
		return
	}

	out, err := h.service.PostPayment(r.Context(), req)
	if err != nil {
		// Basic mapping: validation -> 400, otherwise 502.
		// Store errors are typed; keep it simple without importing store here.
		w.Header().Set("Content-Type", "application/json")
		if isValidationError(err) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(out)
}

func isValidationError(err error) bool {
	// Heuristic: our service returns fmt.Errorf with specific messages.
	// In a larger codebase, prefer typed errors.
	msg := err.Error()
	return msg == "amount must be > 0" || msg == "currency is required" || msg == "order_id is required"
}
