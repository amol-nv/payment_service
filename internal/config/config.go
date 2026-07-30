package config

import "time"

type Config struct {
	PaymentProviderBaseURL string
	PaymentProviderAuthToken string
	HTTPClientTimeout        time.Duration
}

func Default() Config {
	return Config{
		PaymentProviderBaseURL: "http://localhost:8081",
		PaymentProviderAuthToken: "",
		HTTPClientTimeout:        10 * time.Second,
	}
}
