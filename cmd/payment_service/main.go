package main

import (
	"log"
	"net/http"
	"os"

	"amol-nv/payment_service/internal/config"
	"amol-nv/payment_service/internal/mvs/store"
	"amol-nv/payment_service/internal/mvs/view"
	"amol-nv/payment_service/internal/server"
)

func main() {
	cfg := config.Default()
	if v := os.Getenv("PAYMENT_PROVIDER_BASE_URL"); v != "" {
		cfg.PaymentProviderBaseURL = v
	}
	if v := os.Getenv("PAYMENT_PROVIDER_AUTH_TOKEN"); v != "" {
		cfg.PaymentProviderAuthToken = v
	}
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	paymentStore := store.NewPaymentStore(cfg)
	paymentHandler := view.NewPaymentHandler(paymentStore)

	srv := server.New(addr, paymentHandler)
	log.Printf("payment service listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
