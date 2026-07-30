package service

import (
	"context"
	"errors"
	"testing"

	"amol-nv/payment_service/internal/payment/model"
	"amol-nv/payment_service/internal/payment/store"
)

type fakeStore struct {
	fn func(ctx context.Context, req model.PaymentRequest) (model.PaymentResponse, error)
}

func (f *fakeStore) PostPayment(ctx context.Context, req model.PaymentRequest) (model.PaymentResponse, error) {
	return f.fn(ctx, req)
}

func TestService_PostPayment_Validation(t *testing.T) {
	svc := NewService(&fakeStore{fn: func(ctx context.Context, req model.PaymentRequest) (model.PaymentResponse, error) {
		return model.PaymentResponse{}, nil
	}})

	_, err := svc.PostPayment(context.Background(), model.PaymentRequest{Amount: 0, Currency: "USD", Method: "card", OrderID: "ord"})
	if !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestService_PostPayment_ErrorMapping(t *testing.T) {
	svc := NewService(&fakeStore{fn: func(ctx context.Context, req model.PaymentRequest) (model.PaymentResponse, error) {
		return model.PaymentResponse{}, errors.New("wrapped")
	}})

	// Ensure default passthrough when store error isn't recognized
	_, err := svc.PostPayment(context.Background(), model.PaymentRequest{Amount: 1, Currency: "USD", Method: "card", OrderID: "ord"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestService_PostPayment_ErrorMapping_Known(t *testing.T) {
	svc := NewService(&fakeStore{fn: func(ctx context.Context, req model.PaymentRequest) (model.PaymentResponse, error) {
		return model.PaymentResponse{}, store.ErrBadRequest
	}})

	_, err := svc.PostPayment(context.Background(), model.PaymentRequest{Amount: 1, Currency: "USD", Method: "card", OrderID: "ord"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrPaymentBadRequest) {
		t.Fatalf("expected ErrPaymentBadRequest, got %v", err)
	}
}
