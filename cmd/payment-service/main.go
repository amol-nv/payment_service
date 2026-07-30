package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"amol-nv/payment_service/internal/payment/handler"
	"amol-nv/payment_service/internal/payment/mvs"
)

func main() {
	endpointURL := os.Getenv("PAYMENT_ENDPOINT_URL")
	if endpointURL == "" {
		log.Fatal("PAYMENT_ENDPOINT_URL is required")
	}

	authToken := os.Getenv("PAYMENT_AUTH_TOKEN")

	client := &http.Client{Timeout: 10 * time.Second}
	store := mvs.NewPaymentStore(client, endpointURL)
	store.AuthToken = authToken
	store.Timeout = 10 * time.Second

	view := mvs.NewPaymentPostView()
	h := handler.NewPaymentHandler(view, store)

	mux := http.NewServeMux()
	mux.Handle("/payments", h)

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("payment service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
