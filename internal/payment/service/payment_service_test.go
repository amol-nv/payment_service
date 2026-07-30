package service

import (
	"context"
	"errors"
	"testing"

	"amol-nv/payment_service/internal/mvs/dto"
	"amol-nv/payment_service/internal/payment/store"
)

type fakeStore struct {
	called bool
	in      dto.PostPaymentRequest
	out     dto.PostPaymentResponse
	err     error
}

func (f *fakeStore) PostPayment(ctx context.Context, req dto.PostPaymentRequest) (dto.PostPaymentResponse, error) {
	f.called = true
	f.in = req
	return f.out, f.err
}

var _ store.PaymentStore = (*fakeStore)(nil)

func TestPaymentService_PostPayment_Validation(t *testing.T) {
	fs := &fakeStore{}
	svc := NewPaymentService(fs)

	_, err := svc.PostPayment(context.Background(), dto.PostPaymentRequest{Amount: 0, Currency: "USD", OrderID: "o1"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if fs.called {
		t.Fatalf("store should not be called")
	}
}

func TestPaymentService_PostPayment_Orchestrates(t *testing.T) {
	fs := &fakeStore{out: dto.PostPaymentResponse{PaymentID: "p1", Status: "created"}}
	svc := NewPaymentService(fs)

	in := dto.PostPaymentRequest{Amount: 100, Currency: "USD", OrderID: "o1"}
	out, err := svc.PostPayment(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fs.called {
		t.Fatalf("expected store to be called")
	}
	if fs.in != in {
		t.Fatalf("unexpected store input: got %+v want %+v", fs.in, in)
	}
	if out.PaymentID != "p1" {
		t.Fatalf("unexpected output: %+v", out)
	}
}

func TestPaymentService_PostPayment_StoreError(t *testing.T) {
	fs := &fakeStore{err: errors.New("store down")}
	svc := NewPaymentService(fs)

	_, err := svc.PostPayment(context.Background(), dto.PostPaymentRequest{Amount: 100, Currency: "USD", OrderID: "o1"})
	if err == nil {
		t.Fatalf("expected error")
	}
}
