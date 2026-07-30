package store

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"amol-nv/payment_service/internal/payment/model"
)

func TestPaymentStore_PostPayment_Success(t *testing.T) {
	var gotMethod, gotPath, gotAuth, gotContentType string
	var gotBody model.PaymentPostRequest

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")

		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(model.PaymentPostResponse{Status: "ok", TxnID: "txn_123", Message: "accepted"})
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 2 * time.Second}
	st := NewPaymentStore(client, srv.URL, "token123")

	resp, err := st.PostPayment(context.Background(), model.PaymentPostRequest{Amount: 100, Currency: "USD", OrderID: "ord_1", CustomerID: "c1"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.Status != "ok" || resp.TxnID != "txn_123" {
		t.Fatalf("unexpected response: %+v", resp)
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
	if gotContentType != "application/json" {
		t.Fatalf("expected content-type application/json, got %q", gotContentType)
	}
	if gotBody.Amount != 100 || gotBody.Currency != "USD" || gotBody.OrderID != "ord_1" || gotBody.CustomerID != "c1" {
		t.Fatalf("unexpected request body: %+v", gotBody)
	}
}

func TestPaymentStore_PostPayment_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad_request","message":"invalid"}`))
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 2 * time.Second}
	st := NewPaymentStore(client, srv.URL, "")

	_, err := st.PostPayment(context.Background(), model.PaymentPostRequest{Amount: 1, Currency: "USD", OrderID: "o"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "payment post failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}
