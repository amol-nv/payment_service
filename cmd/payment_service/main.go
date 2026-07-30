package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"amol-nv/payment_service/internal/payment/service"
	"amol-nv/payment_service/internal/payment/store"
	"amol-nv/payment_service/internal/payment/view"
)

func main() {
	providerBaseURL := os.Getenv("PAYMENT_PROVIDER_BASE_URL")
	if providerBaseURL == "" {
		providerBaseURL = "http://localhost:8081"
	}
	providerEndpointPath := os.Getenv("PAYMENT_PROVIDER_ENDPOINT_PATH")
	if providerEndpointPath == "" {
		providerEndpointPath = "/payments"
	}

	st, err := store.NewHTTPPaymentStore(store.PaymentStoreConfig{
		BaseURL:      providerBaseURL,
		EndpointPath: providerEndpointPath,
		Timeout:      10 * time.Second,
	})
	if err != nil {
		log.Fatalf("failed to create payment store: %v", err)
	}

	svc := service.NewPaymentService(st)
	h := view.NewPaymentHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/payments", h.PostPayment)

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("payment service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
