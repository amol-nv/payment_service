package mvs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPDoer is the minimal interface needed from *http.Client.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Store performs the external POST call.
// It is responsible for building the HTTP request, executing it, and returning typed response data.
type PaymentStore struct {
	Client        HTTPDoer
	EndpointURL  string
	AuthToken    string
	Timeout       time.Duration
	ContentType   string
	UserAgent     string
	AdditionalHdr map[string]string
}

func NewPaymentStore(client HTTPDoer, endpointURL string) *PaymentStore {
	return &PaymentStore{
		Client:       client,
		EndpointURL: endpointURL,
		Timeout:      10 * time.Second,
		ContentType:  "application/json",
	}
}

// Post executes the outbound POST call using the MVS Model.
func (s *PaymentStore) Post(ctx context.Context, reqModel PaymentPostRequest) (PaymentPostResponse, error) {
	if s.Client == nil {
		return PaymentPostResponse{}, errors.New("payment store: http client is nil")
	}
	if s.EndpointURL == "" {
		return PaymentPostResponse{}, errors.New("payment store: endpoint url is empty")
	}

	bodyBytes, err := json.Marshal(reqModel)
	if err != nil {
		return PaymentPostResponse{}, fmt.Errorf("payment store: marshal request: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.EndpointURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return PaymentPostResponse{}, fmt.Errorf("payment store: create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", s.ContentType)
	if s.UserAgent != "" {
		httpReq.Header.Set("User-Agent", s.UserAgent)
	}
	if s.AuthToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+s.AuthToken)
	}
	for k, v := range s.AdditionalHdr {
		if k == "" {
			continue
		}
		httpReq.Header.Set(k, v)
	}

	resp, err := s.Client.Do(httpReq)
	if err != nil {
		return PaymentPostResponse{}, fmt.Errorf("payment store: do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return PaymentPostResponse{}, fmt.Errorf("payment store: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Best-effort parse error payload; otherwise include raw body.
		var parsed map[string]any
		if err := json.Unmarshal(respBody, &parsed); err == nil {
			return PaymentPostResponse{}, fmt.Errorf("payment store: non-2xx response: status=%d body=%v", resp.StatusCode, parsed)
		}
		return PaymentPostResponse{}, fmt.Errorf("payment store: non-2xx response: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var out PaymentPostResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return PaymentPostResponse{}, fmt.Errorf("payment store: unmarshal response: %w", err)
	}
	return out, nil
}
