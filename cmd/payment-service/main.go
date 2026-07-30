package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"amol-nv/payment_service/internal/config"
	"amol-nv/payment_service/internal/httpapi"
	"amol-nv/payment_service/internal/payment"
	"amol-nv/payment_service/internal/service"
)

func main() {
	cfg := config.FromEnv()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	paymentService := service.NewPaymentService(client, cfg)
	handler := httpapi.NewPaymentHandler(paymentService)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/payments", handler.PostPayment)

	addr := cfg.ListenAddr
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("payment service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
	_ = os.Stdout
}
