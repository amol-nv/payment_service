package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"amol-nv/payment_service/internal/config"
	"amol-nv/payment_service/internal/payment"
)

type PaymentService interface {
	PostPayment(ctx context.Context, req payment.PostPaymentRequest) (payment.PostPaymentResponse, error)
}

type paymentService struct {
	client *http.Client
	cfg    config.Config
}

func NewPaymentService(client *http.Client, cfg config.Config) PaymentService {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &paymentService{client: client, cfg: cfg}
}

func (s *paymentService) PostPayment(ctx context.Context, req payment.PostPaymentRequest) (payment.PostPaymentResponse, error) {
	if strings.TrimSpace(s.cfg.UpstreamBaseURL) == "" {
		return payment.PostPaymentResponse{}, errors.New("PAYMENT_UPSTREAM_BASE_URL is not configured")
	}

	base, err := url.Parse(strings.TrimRight(s.cfg.UpstreamBaseURL, "/"))
	if err != nil {
		return payment.PostPaymentResponse{}, fmt.Errorf("invalid upstream base url: %w", err)
	}
	path := s.cfg.UpstreamPath
	if path == "" {
		path = "/payments"
	}
	upstreamURL := base.ResolveReference(&url.URL{Path: path})

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return payment.PostPaymentResponse{}, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return payment.PostPaymentResponse{}, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(s.cfg.AuthToken) != "" {
		httpReq.Header.Set("Authorization", "Bearer "+s.cfg.AuthToken)
	}

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return payment.PostPaymentResponse{}, fmt.Errorf("upstream post failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return payment.PostPaymentResponse{}, fmt.Errorf("read upstream response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return payment.PostPaymentResponse{}, fmt.Errorf("upstream returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var out payment.PostPaymentResponse
	if err := json.Unmarshal(respBytes, &out); err != nil {
		// If upstream returns a different schema, still surface raw payload.
		out.Raw = json.RawMessage(respBytes)
		return out, fmt.Errorf("unmarshal upstream response: %w", err)
	}
	out.Raw = json.RawMessage(respBytes)
	return out, nil
}
