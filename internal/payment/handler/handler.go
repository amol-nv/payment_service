package handler

import (
	"context"
	"net/http"
	"time"

	"amol-nv/payment_service/internal/payment/mvs"
)

// PaymentHandler wires the MVS flow into an HTTP handler.
type PaymentHandler struct {
	View  *mvs.PaymentPostView
	Store *mvs.PaymentStore
	// RequestTimeout bounds the overall handler execution.
	RequestTimeout time.Duration
}

func NewPaymentHandler(view *mvs.PaymentPostView, store *mvs.PaymentStore) *PaymentHandler {
	return &PaymentHandler{
		View:           view,
		Store:          store,
		RequestTimeout: 15 * time.Second,
	}
}

// ServeHTTP handles POST /payments.
// It maps incoming JSON -> Model (View), calls Store.Post, then maps Model -> JSON response (View).
func (h *PaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if h.View == nil || h.Store == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	if h.RequestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.RequestTimeout)
		defer cancel()
	}

	payload, err := readAllLimit(r, 1<<20) // 1MB
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	modelReq, err := h.View.ToModel(payload)
	if err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	modelResp, err := h.Store.Post(ctx, modelReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	respBytes, err := h.View.FromModel(modelResp)
	if err != nil {
		http.Error(w, "failed to build response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respBytes)
}
