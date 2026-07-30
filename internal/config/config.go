package config

import "os"

type Config struct {
	ListenAddr string

	UpstreamBaseURL string
	UpstreamPath    string
	AuthToken       string
}

func FromEnv() Config {
	return Config{
		ListenAddr:      getenv("PAYMENT_LISTEN_ADDR", ":8080"),
		UpstreamBaseURL: getenv("PAYMENT_UPSTREAM_BASE_URL", ""),
		UpstreamPath:    getenv("PAYMENT_UPSTREAM_PATH", "/payments"),
		AuthToken:       os.Getenv("PAYMENT_UPSTREAM_AUTH_TOKEN"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
