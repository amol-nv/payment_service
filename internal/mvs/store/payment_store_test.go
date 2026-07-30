package store

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"amol-nv/payment_service/internal/config"
	"amol-nv/payment_service/internal/mvs/model"
)

func TestPaymentStore_CreatePayment_Success(t *testing.T) {
	var gotMethod, gotPath, gotAuth string
	var gotBody model.PaymentPostRequest

	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")

		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(model.PaymentPostResponse{
			PaymentID: "pay_123",
			Status:    "created",
			Message:   "ok",
		})
	}))
	defer provider.Close()

	cfg := config.Default()
	cfg.PaymentProviderBaseURL = provider.URL
	cfg.PaymentProviderAuthToken = "token123"

	s := NewPaymentStore(cfg)
	resp, status, err := s.CreatePayment(context.Background(), model.PaymentPostRequest{
		CustomerID: "c1",
		Amount:     10.5,
		Currency:   "USD",
		Reference:  "ref",
	})
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d", status)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("expected method POST, got %s", gotMethod)
	}
	if gotPath != "/payments" {
		t.Fatalf("expected path /payments, got %s", gotPath)
	}
	if gotAuth != "Bearer token123" {
		t.Fatalf("expected auth header, got %q", gotAuth)
	}
	if resp.PaymentID != "pay_123" || resp.Status != "created" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if gotBody.CustomerID != "c1" || gotBody.Currency != "USD" || gotBody.Reference != "ref" {
		t.Fatalf("unexpected request body: %+v", gotBody)
	}
}

func TestPaymentStore_CreatePayment_ProviderError(t *testing.T) {
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
	}))
	defer provider.Close()

	cfg := config.Default()
	cfg.PaymentProviderBaseURL = provider.URL

	s := NewPaymentStore(cfg)
	_, status, err := s.CreatePayment(context.Background(), model.PaymentPostRequest{
		CustomerID: "c1",
		Amount:     1,
		Currency:   "USD",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", status)
	}
}
