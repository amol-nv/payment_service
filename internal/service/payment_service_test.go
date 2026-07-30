package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"amol-nv/payment_service/internal/config"
	"amol-nv/payment_service/internal/payment"
)

func TestPaymentService_PostPayment_SendsPostAndParsesResponse(t *testing.T) {
	var gotMethod string
	var gotAuth string
	var gotBody payment.PostPaymentRequest

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotAuth = r.Header.Get("Authorization")
		_ = json.NewDecoder(r.Body).Decode(&gotBody)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payment.PostPaymentResponse{
			PaymentID: "pay_123",
			Status:    "SUCCESS",
		})
	}))
	defer upstream.Close()

	cfg := config.Config{
		UpstreamBaseURL: upstream.URL,
		UpstreamPath:    "/payments",
		AuthToken:       "token123",
	}

	client := &http.Client{Timeout: 2 * time.Second}
	svc := NewPaymentService(client, cfg)

	req := payment.PostPaymentRequest{CustomerID: "c1", Amount: 10.5, Currency: "USD", Reference: "r1"}
	resp, err := svc.PostPayment(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("expected method POST, got %s", gotMethod)
	}
	if gotAuth != "Bearer token123" {
		t.Fatalf("expected auth header, got %q", gotAuth)
	}
	if gotBody.CustomerID != req.CustomerID || gotBody.Amount != req.Amount || gotBody.Currency != req.Currency || gotBody.Reference != req.Reference {
		t.Fatalf("unexpected body: %+v", gotBody)
	}
	if resp.PaymentID != "pay_123" || resp.Status != "SUCCESS" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestPaymentService_PostPayment_UpstreamError(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad"}`))
	}))
	defer upstream.Close()

	cfg := config.Config{UpstreamBaseURL: upstream.URL, UpstreamPath: "/payments"}
	svc := NewPaymentService(&http.Client{Timeout: 2 * time.Second}, cfg)

	_, err := svc.PostPayment(context.Background(), payment.PostPaymentRequest{CustomerID: "c1", Amount: 1, Currency: "USD"})
	if err == nil {
		t.Fatalf("expected error")
	}
}
