package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"amol-nv/payment_service/internal/payment/model"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Store struct {
	client  HTTPClient
	baseURL string
	apiKey  string
	// endpoint path for payment post
	paymentPath string
	timeout     time.Duration
}

type StoreOption func(*Store)

func WithPaymentPath(path string) StoreOption {
	return func(s *Store) { s.paymentPath = path }
}

func WithTimeout(d time.Duration) StoreOption {
	return func(s *Store) { s.timeout = d }
}

func NewStore(client HTTPClient, baseURL, apiKey string, opts ...StoreOption) *Store {
	s := &Store{
		client:      client,
		baseURL:     baseURL,
		apiKey:      apiKey,
		paymentPath: "/payments",
		timeout:     10 * time.Second,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type PaymentStore interface {
	PostPayment(ctx context.Context, req model.PaymentRequest) (model.PaymentResponse, error)
}

func (s *Store) PostPayment(ctx context.Context, req model.PaymentRequest) (model.PaymentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	url := fmt.Sprintf("%s%s", s.baseURL, s.paymentPath)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return model.PaymentResponse{}, fmt.Errorf("marshal payment request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return model.PaymentResponse{}, fmt.Errorf("create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	}
	if req.IdempotencyKey != "" {
		httpReq.Header.Set("Idempotency-Key", req.IdempotencyKey)
	}

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return model.PaymentResponse{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return model.PaymentResponse{}, fmt.Errorf("read response body: %w", readErr)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to decode structured error
		var er model.ErrorResponse
		if len(respBody) > 0 {
			_ = json.Unmarshal(respBody, &er)
		}

		switch resp.StatusCode {
		case http.StatusBadRequest:
			return model.PaymentResponse{}, fmt.Errorf("%w: %s", ErrBadRequest, er.Message)
		case http.StatusUnauthorized:
			return model.PaymentResponse{}, fmt.Errorf("%w: %s", ErrUnauthorized, er.Message)
		case http.StatusForbidden:
			return model.PaymentResponse{}, fmt.Errorf("%w: %s", ErrForbidden, er.Message)
		case http.StatusNotFound:
			return model.PaymentResponse{}, fmt.Errorf("%w: %s", ErrNotFound, er.Message)
		case http.StatusConflict:
			return model.PaymentResponse{}, fmt.Errorf("%w: %s", ErrConflict, er.Message)
		case http.StatusUnprocessableEntity:
			return model.PaymentResponse{}, fmt.Errorf("%w: %s", ErrUnprocessable, er.Message)
		default:
			return model.PaymentResponse{}, fmt.Errorf("%w (status=%d): %s", ErrUpstreamFailure, resp.StatusCode, string(respBody))
		}
	}

	var out model.PaymentResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return model.PaymentResponse{}, fmt.Errorf("decode payment response: %w", err)
	}
	return out, nil
}
