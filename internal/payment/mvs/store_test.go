package mvs

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }

func TestPaymentStore_Post_Success(t *testing.T) {
	var gotMethod, gotAuth, gotContentType string
	var gotBody string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(PaymentPostResponse{Status: "ok", PaymentID: "pmt_123", Message: "accepted"})
	}))
	defer srv.Close()

	store := NewPaymentStore(srv.Client(), srv.URL)
	store.AuthToken = "token123"

	resp, err := store.Post(context.Background(), PaymentPostRequest{Amount: 100, Currency: "USD", OrderID: "ord_1"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("expected method POST, got %s", gotMethod)
	}
	if gotAuth != "Bearer token123" {
		t.Fatalf("expected auth header, got %q", gotAuth)
	}
	if gotContentType != "application/json" {
		t.Fatalf("expected content-type application/json, got %q", gotContentType)
	}
	if !strings.Contains(gotBody, "\"amount\":100") {
		t.Fatalf("expected request body to contain amount, got %s", gotBody)
	}
	if resp.Status != "ok" || resp.PaymentID != "pmt_123" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestPaymentStore_Post_Non2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer srv.Close()

	store := NewPaymentStore(srv.Client(), srv.URL)
	_, err := store.Post(context.Background(), PaymentPostRequest{Amount: 1, Currency: "USD", OrderID: "o"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "non-2xx") {
		t.Fatalf("expected non-2xx error, got %v", err)
	}
}
