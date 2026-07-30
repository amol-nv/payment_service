package store

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"amol-nv/payment_service/internal/mvs/dto"
)

func TestHTTPPaymentStore_PostPayment_Success(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody []byte

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotBody, _ = io.ReadAll(r.Body)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(dto.PostPaymentResponse{PaymentID: "p1", Status: "created"})
	}))
	defer ts.Close()

	st, err := NewHTTPPaymentStore(PaymentStoreConfig{
		BaseURL:      ts.URL,
		EndpointPath: "/payments",
		Timeout:      2 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewHTTPPaymentStore error: %v", err)
	}

	in := dto.PostPaymentRequest{Amount: 100, Currency: "USD", OrderID: "o1"}
	out, err := st.PostPayment(context.Background(), in)
	if err != nil {
		t.Fatalf("PostPayment error: %v", err)
	}
	if out.PaymentID != "p1" || out.Status != "created" {
		t.Fatalf("unexpected response: %+v", out)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("expected method POST, got %s", gotMethod)
	}
	if gotPath != "/payments" {
		t.Fatalf("expected path /payments, got %s", gotPath)
	}

	var decoded dto.PostPaymentRequest
	if err := json.Unmarshal(gotBody, &decoded); err != nil {
		t.Fatalf("invalid json body: %v", err)
	}
	if decoded != in {
		t.Fatalf("unexpected request body: got %+v want %+v", decoded, in)
	}
}

func TestHTTPPaymentStore_PostPayment_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Code: "bad_request", Message: "bad"})
	}))
	defer ts.Close()

	st, err := NewHTTPPaymentStore(PaymentStoreConfig{
		BaseURL:      ts.URL,
		EndpointPath: "/payments",
		Timeout:      2 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewHTTPPaymentStore error: %v", err)
	}

	_, err = st.PostPayment(context.Background(), dto.PostPaymentRequest{Amount: 1, Currency: "USD", OrderID: "o1"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if _, ok := err.(*HTTPError); !ok {
		t.Fatalf("expected *HTTPError, got %T", err)
	}
}
