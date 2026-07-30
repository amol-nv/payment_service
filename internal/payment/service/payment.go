package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"amol-nv/payment_service/internal/payment/model"
	"amol-nv/payment_service/internal/payment/store"
)

type PaymentService interface {
	PostPayment(ctx context.Context, req model.PaymentRequest) (model.PaymentResponse, error)
}

type Service struct {
	store store.PaymentStore
}

func NewService(store store.PaymentStore) *Service {
	return &Service{store: store}
}

func (s *Service) PostPayment(ctx context.Context, req model.PaymentRequest) (model.PaymentResponse, error) {
	if req.Amount <= 0 {
		return model.PaymentResponse{}, ErrInvalidAmount
	}
	if strings.TrimSpace(req.Currency) == "" {
		return model.PaymentResponse{}, ErrInvalidCurrency
	}
	if strings.TrimSpace(req.Method) == "" {
		return model.PaymentResponse{}, ErrInvalidMethod
	}
	if strings.TrimSpace(req.OrderID) == "" {
		return model.PaymentResponse{}, ErrInvalidOrderID
	}

	out, err := s.store.PostPayment(ctx, req)
	if err == nil {
		return out, nil
	}

	// Map store errors to service errors
	switch {
	case errors.Is(err, store.ErrBadRequest):
		return model.PaymentResponse{}, fmt.Errorf("%w: %v", ErrPaymentBadRequest, err)
	case errors.Is(err, store.ErrUnauthorized):
		return model.PaymentResponse{}, fmt.Errorf("%w: %v", ErrPaymentUnauthorized, err)
	case errors.Is(err, store.ErrForbidden):
		return model.PaymentResponse{}, fmt.Errorf("%w: %v", ErrPaymentForbidden, err)
	case errors.Is(err, store.ErrNotFound):
		return model.PaymentResponse{}, fmt.Errorf("%w: %v", ErrPaymentNotFound, err)
	case errors.Is(err, store.ErrConflict):
		return model.PaymentResponse{}, fmt.Errorf("%w: %v", ErrPaymentConflict, err)
	case errors.Is(err, store.ErrUnprocessable):
		return model.PaymentResponse{}, fmt.Errorf("%w: %v", ErrPaymentUnprocessable, err)
	case errors.Is(err, store.ErrUpstreamFailure):
		return model.PaymentResponse{}, fmt.Errorf("%w: %v", ErrPaymentUpstreamFailure, err)
	default:
		return model.PaymentResponse{}, err
	}
}
