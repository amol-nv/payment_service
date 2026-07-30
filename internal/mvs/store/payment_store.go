package store

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"amol-nv/payment_service/internal/config"
	"amol-nv/payment_service/internal/mvs/model"
)

type PaymentStore interface {
	CreatePayment(ctx context.Context, req model.PaymentPostRequest) (model.PaymentPostResponse, int, error)
}

type paymentStore struct {
	client  *http.Client
	baseURL string
	authTok string
}

func NewPaymentStore(cfg config.Config) PaymentStore {
	client := &http.Client{Timeout: cfg.HTTPClientTimeout}
	return &paymentStore{
		client:  client,
		baseURL: strings.TrimRight(cfg.PaymentProviderBaseURL, "/"),
		authTok: cfg.PaymentProviderAuthToken,
	}
}

func (s *paymentStore) CreatePayment(ctx context.Context, req model.PaymentPostRequest) (model.PaymentPostResponse, int, error) {
	var out model.PaymentPostResponse

	b, err := json.Marshal(req)
	if err != nil {
		return out, http.StatusInternalServerError, fmt.Errorf("marshal request: %w", err)
	}

	url := s.baseURL + "/payments"
	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return out, http.StatusInternalServerError, fmt.Errorf("create request: %w", err)
	}
	hreq.Header.Set("Content-Type", "application/json")
	if s.authTok != "" {
		hreq.Header.Set("Authorization", "Bearer "+s.authTok)
	}

	resp, err := s.client.Do(hreq)
	if err != nil {
		// Map transport errors to 502
		return out, http.StatusBadGateway, fmt.Errorf("payment provider call failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Best-effort parse error message
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return out, resp.StatusCode, errors.New(msg)
	}

	if err := json.Unmarshal(body, &out); err != nil {
		return out, http.StatusBadGateway, fmt.Errorf("unmarshal response: %w", err)
	}

	return out, resp.StatusCode, nil
}

// Ensure we reference time to avoid unused import if config changes.
var _ = time.Second
