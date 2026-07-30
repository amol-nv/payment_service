package service

import (
	"context"
	"fmt"

	"amol-nv/payment_service/internal/mvs/dto"
	"amol-nv/payment_service/internal/payment/store"
)

type PaymentService interface {
	PostPayment(ctx context.Context, req dto.PostPaymentRequest) (dto.PostPaymentResponse, error)
}

type paymentService struct {
	store store.PaymentStore
}

func NewPaymentService(st store.PaymentStore) PaymentService {
	return &paymentService{store: st}
}

func (s *paymentService) PostPayment(ctx context.Context, req dto.PostPaymentRequest) (dto.PostPaymentResponse, error) {
	if req.Amount <= 0 {
		return dto.PostPaymentResponse{}, fmt.Errorf("amount must be > 0")
	}
	if req.Currency == "" {
		return dto.PostPaymentResponse{}, fmt.Errorf("currency is required")
	}
	if req.OrderID == "" {
		return dto.PostPaymentResponse{}, fmt.Errorf("order_id is required")
	}

	return s.store.PostPayment(ctx, req)
}
