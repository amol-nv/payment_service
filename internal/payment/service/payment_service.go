package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"amol-nv/payment_service/internal/payment/model"
	"amol-nv/payment_service/internal/payment/store"
)

type PaymentService struct {
	store store.PaymentStore
}

func NewPaymentService(st store.PaymentStore) *PaymentService {
	return &PaymentService{store: st}
}

func (s *PaymentService) HandlePostPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req model.PaymentPostRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{Error: "invalid_request", Message: err.Error()})
		return
	}

	if err := validate(req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{Error: "validation_error", Message: err.Error()})
		return
	}

	resp, err := s.store.PostPayment(r.Context(), req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		// Keep mapping minimal; treat store errors as bad gateway.
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{Error: "payment_post_failed", Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func validate(req model.PaymentPostRequest) error {
	if req.Amount <= 0 {
		return errors.New("amount must be > 0")
	}
	if strings.TrimSpace(req.Currency) == "" {
		return errors.New("currency is required")
	}
	if strings.TrimSpace(req.OrderID) == "" {
		return errors.New("order_id is required")
	}
	return nil
}

func (s *PaymentService) PostPayment(ctx context.Context, req model.PaymentPostRequest) (model.PaymentPostResponse, error) {
	if err := validate(req); err != nil {
		return model.PaymentPostResponse{}, fmt.Errorf("%w", err)
	}
	return s.store.PostPayment(ctx, req)
}
