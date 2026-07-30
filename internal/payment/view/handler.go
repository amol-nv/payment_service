package view

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"amol-nv/payment_service/internal/payment/model"
	"amol-nv/payment_service/internal/payment/service"
)

type Handler struct {
	service service.PaymentService
}

func NewHandler(s service.PaymentService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) PostPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req model.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{Code: "bad_request", Message: "invalid json"})
		return
	}

	out, err := h.service.PostPayment(context.Background(), req)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrInvalidAmount), errors.Is(err, service.ErrInvalidCurrency), errors.Is(err, service.ErrInvalidMethod), errors.Is(err, service.ErrInvalidOrderID):
			status = http.StatusBadRequest
		case errors.Is(err, service.ErrPaymentBadRequest):
			status = http.StatusBadRequest
		case errors.Is(err, service.ErrPaymentUnauthorized):
			status = http.StatusUnauthorized
		case errors.Is(err, service.ErrPaymentForbidden):
			status = http.StatusForbidden
		case errors.Is(err, service.ErrPaymentNotFound):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrPaymentConflict):
			status = http.StatusConflict
		case errors.Is(err, service.ErrPaymentUnprocessable):
			status = http.StatusUnprocessableEntity
		case errors.Is(err, service.ErrPaymentUpstreamFailure):
			status = http.StatusBadGateway
		}

		writeJSON(w, status, model.ErrorResponse{Code: "payment_error", Message: err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, out)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
