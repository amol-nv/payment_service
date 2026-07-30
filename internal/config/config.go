package config

import (
	"os"
	"time"
)

type Config struct {
	PaymentBaseURL string
	HTTPTimeout    time.Duration
	AuthToken      string
}

func Load() Config {
	baseURL := os.Getenv("PAYMENT_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}

	timeout := 10 * time.Second
	if v := os.Getenv("PAYMENT_HTTP_TIMEOUT_SECONDS"); v != "" {
		// best-effort parse; keep default on error
		if parsed, err := time.ParseDuration(v + "s"); err == nil {
			timeout = parsed
		}
	}

	return Config{
		PaymentBaseURL: baseURL,
		HTTPTimeout:    timeout,
		AuthToken:      os.Getenv("PAYMENT_AUTH_TOKEN"),
	}
}
