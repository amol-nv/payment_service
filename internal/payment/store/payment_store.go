package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"amol-nv/payment_service/internal/mvs/dto"
)

type PaymentStore interface {
	PostPayment(ctx context.Context, req dto.PostPaymentRequest) (dto.PostPaymentResponse, error)
}

type httpPaymentStore struct {
	client  *http.Client
	baseURL string
	// endpointPath is the path on the payment provider.
	endpointPath string
}

type PaymentStoreConfig struct {
	BaseURL      string
	EndpointPath string
	Timeout      time.Duration
}

func NewHTTPPaymentStore(cfg PaymentStoreConfig) (PaymentStore, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("BaseURL is required")
	}
	if cfg.EndpointPath == "" {
		cfg.EndpointPath = "/payments"
	}
	client := &http.Client{Timeout: cfg.Timeout}
	return &httpPaymentStore{
		client:       client,
		baseURL:      cfg.BaseURL,
		endpointPath: cfg.EndpointPath,
	}, nil
}

func (s *httpPaymentStore) PostPayment(ctx context.Context, req dto.PostPaymentRequest) (dto.PostPaymentResponse, error) {
	url := s.baseURL + s.endpointPath

	b, err := json.Marshal(req)
	if err != nil {
		return dto.PostPaymentResponse{}, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return dto.PostPaymentResponse{}, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return dto.PostPaymentResponse{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.PostPaymentResponse{}, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to decode provider error.
		var er dto.ErrorResponse
		if len(body) > 0 {
			_ = json.Unmarshal(body, &er)
		}
		msg := er.Message
		return dto.PostPaymentResponse{}, &HTTPError{StatusCode: resp.StatusCode, Body: body, Message: msg}
	}

	var out dto.PostPaymentResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return dto.PostPaymentResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return out, nil
}
