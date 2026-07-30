package main

import (
	"log"
	"net/http"
	"os"

	"amol-nv/payment_service/internal/config"
	"amol-nv/payment_service/internal/httpclient"
	"amol-nv/payment_service/internal/payment/router"
	"amol-nv/payment_service/internal/payment/service"
	"amol-nv/payment_service/internal/payment/store"
)

func main() {
	cfg := config.Load()
	client := httpclient.NewClient(cfg.HTTPTimeout)
	st := store.NewPaymentStore(client, cfg.PaymentBaseURL, cfg.AuthToken)
	svc := service.NewPaymentService(st)
	r := router.NewRouter(svc)

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}

	log.Printf("payment_service listening on :%s", addr)
	if err := http.ListenAndServe(":"+addr, r.Handler()); err != nil {
		log.Fatal(err)
	}
}
