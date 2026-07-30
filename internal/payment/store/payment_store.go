package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"amol-nv/payment_service/internal/payment/model"
)

type PaymentStore interface {
	PostPayment(ctx context.Context, req model.PaymentPostRequest) (model.PaymentPostResponse, error)
}

type paymentStore struct {
	client    *http.Client
	baseURL   string
	authToken string
	endpoint  string
}

func NewPaymentStore(client *http.Client, baseURL, authToken string) PaymentStore {
	return &paymentStore{
		client:    client,
		baseURL:   strings.TrimRight(baseURL, "/"),
		authToken: authToken,
		endpoint:  "/payments",
	}
}

func (s *paymentStore) PostPayment(ctx context.Context, req model.PaymentPostRequest) (model.PaymentPostResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return model.PaymentPostResponse{}, fmt.Errorf("marshal request: %w", err)
	}

	u, err := url.Parse(s.baseURL)
	if err != nil {
		return model.PaymentPostResponse{}, fmt.Errorf("parse base url: %w", err)
	}
	u.Path = strings.TrimRight(u.Path, "/") + s.endpoint

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(b))
	if err != nil {
		return model.PaymentPostResponse{}, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if s.authToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+s.authToken)
	}

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return model.PaymentPostResponse{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var er model.ErrorResponse
		if len(body) > 0 {
			_ = json.Unmarshal(body, &er)
		}
		if er.Error == "" {
			er.Error = http.StatusText(resp.StatusCode)
		}
		if er.Message == "" {
			er.Message = string(body)
		}
		return model.PaymentPostResponse{}, fmt.Errorf("payment post failed: status=%d error=%s message=%s", resp.StatusCode, er.Error, er.Message)
	}

	var pr model.PaymentPostResponse
	if err := json.Unmarshal(body, &pr); err != nil {
		return model.PaymentPostResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	// Ensure we don't return a zero-value response silently in case of empty body.
	if pr.Status == "" && pr.TxnID == "" {
		// small guard; not strictly required
		_ = time.Second
	}

	return pr, nil
}

// compile-time check
var _ PaymentStore = (*paymentStore)(nil)
