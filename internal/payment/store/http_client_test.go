package store

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"amol-nv/payment_service/internal/payment/model"
)

type recordingClient struct {
	baseURL string
	lastReq *http.Request
}

func (c *recordingClient) Do(req *http.Request) (*http.Response, error) {
	c.lastReq = req
	// Use default transport to hit the httptest server
	return http.DefaultClient.Do(req)
}

func TestStore_PostPayment_Success(t *testing.T) {
	var gotAuth, gotIdempotency, gotContentType string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotIdempotency = r.Header.Get("Idempotency-Key")
		gotContentType = r.Header.Get("Content-Type")

		b, _ := io.ReadAll(r.Body)
		var pr model.PaymentRequest
		if err := json.Unmarshal(b, &pr); err != nil {
			t.Fatalf("invalid json: %v", err)
		}
		if pr.OrderID != "ord-1" {
			t.Fatalf("unexpected order id: %s", pr.OrderID)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(model.PaymentResponse{
			PaymentID: "pay-1",
			Status:    "created",
			CreatedAt: time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC),
		})
	}))
	defer srv.Close()

	client := &recordingClient{}
	st := NewStore(client, srv.URL, "test-api-key", WithTimeout(2*time.Second))

	out, err := st.PostPayment(context.Background(), model.PaymentRequest{
		Amount:          100,
		Currency:        "USD",
		Method:          "card",
		OrderID:         "ord-1",
		IdempotencyKey: "idem-1",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if out.PaymentID != "pay-1" {
		t.Fatalf("unexpected payment id: %s", out.PaymentID)
	}
	if gotAuth != "Bearer test-api-key" {
		t.Fatalf("unexpected auth header: %s", gotAuth)
	}
	if gotIdempotency != "idem-1" {
		t.Fatalf("unexpected idempotency header: %s", gotIdempotency)
	}
	if gotContentType != "application/json" {
		t.Fatalf("unexpected content-type: %s", gotContentType)
	}
}

func TestStore_PostPayment_ErrorMapping(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{Code: "invalid", Message: "bad payload"})
	}))
	defer srv.Close()

	st := NewStore(http.DefaultClient, srv.URL, "", WithTimeout(2*time.Second))
	_, err := st.PostPayment(context.Background(), model.PaymentRequest{Amount: 1, Currency: "USD", Method: "card", OrderID: "ord"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errorsIs(err, ErrUnprocessable) {
		t.Fatalf("expected ErrUnprocessable, got %v", err)
	}
}

// small helper to avoid importing errors in test file twice
func errorsIs(err, target error) bool {
	return err != nil && (err == target || (interface{ Unwrap() error }(err) != nil && errorsIs(err.(interface{ Unwrap() error }).Unwrap(), target)))
}
