package view

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"amol-nv/payment_service/internal/mvs/model"
	"amol-nv/payment_service/internal/mvs/store"
)

type fakeStore struct {
	createFn func(ctx context.Context, req model.PaymentPostRequest) (model.PaymentPostResponse, int, error)
}

func (f *fakeStore) CreatePayment(ctx context.Context, req model.PaymentPostRequest) (model.PaymentPostResponse, int, error) {
	return f.createFn(ctx, req)
}

func TestPaymentHandler_CreatePayment_Success(t *testing.T) {
	fs := &fakeStore{createFn: func(ctx context.Context, req model.PaymentPostRequest) (model.PaymentPostResponse, int, error) {
		return model.PaymentPostResponse{PaymentID: "p1", Status: "created"}, http.StatusOK, nil
	}}
	h := NewPaymentHandler(fs)

	body, _ := json.Marshal(model.PaymentPostRequest{CustomerID: "c1", Amount: 2, Currency: "USD"})
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreatePayment(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var got model.PaymentPostResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.PaymentID != "p1" || got.Status != "created" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestPaymentHandler_CreatePayment_InvalidJSON(t *testing.T) {
	fs := &fakeStore{createFn: func(ctx context.Context, req model.PaymentPostRequest) (model.PaymentPostResponse, int, error) {
		return model.PaymentPostResponse{}, http.StatusInternalServerError, nil
	}}
	h := NewPaymentHandler(fs)

	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	h.CreatePayment(w, req)
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestPaymentHandler_CreatePayment_StoreError(t *testing.T) {
	fs := &fakeStore{createFn: func(ctx context.Context, req model.PaymentPostRequest) (model.PaymentPostResponse, int, error) {
		return model.PaymentPostResponse{}, http.StatusBadGateway, http.ErrHandlerTimeout
	}}
	h := NewPaymentHandler(fs)

	body, _ := json.Marshal(model.PaymentPostRequest{CustomerID: "c1", Amount: 2, Currency: "USD"})
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreatePayment(w, req)
	if w.Result().StatusCode != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", w.Result().StatusCode)
	}
}

// compile-time check
var _ store.PaymentStore = (*fakeStore)(nil)
