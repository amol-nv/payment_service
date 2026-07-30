package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"amol-nv/payment_service/internal/payment/mvs"
)

type fakeStore struct {
	gotReq mvs.PaymentPostRequest
	resp   mvs.PaymentPostResponse
	err    error
}

func (f *fakeStore) Post(ctx context.Context, reqModel mvs.PaymentPostRequest) (mvs.PaymentPostResponse, error) {
	f.gotReq = reqModel
	return f.resp, f.err
}

func TestPaymentHandler_Success(t *testing.T) {
	view := mvs.NewPaymentPostView()
	fs := &fakeStore{resp: mvs.PaymentPostResponse{Status: "ok", PaymentID: "p_1"}}

	// Wrap fakeStore into PaymentStore-like struct by using a real PaymentStore is not possible.
	// Instead, we create a minimal adapter by embedding into a struct with same method set.
	// We'll use a small trick: create a PaymentStore and override its Client with a fake RoundTripper.
	// However, to keep this test focused, we will use httptest server via real PaymentStore.
	//
	// For simplicity, use httptest server.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &fs.gotReq)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(fs.resp)
	}))
	defer srv.Close()

	store := mvs.NewPaymentStore(srv.Client(), srv.URL)
	h := NewPaymentHandler(view, store)

	payload := []byte(`{"amount":10,"currency":"USD","order_id":"ord_1"}`)
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewReader(payload))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	if fs.gotReq.Amount != 10 || fs.gotReq.Currency != "USD" || fs.gotReq.OrderID != "ord_1" {
		t.Fatalf("unexpected got request: %+v", fs.gotReq)
	}
}
