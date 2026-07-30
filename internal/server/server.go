package server

import (
	"net/http"

	"amol-nv/payment_service/internal/mvs/view"
)

type Server struct {
	httpServer *http.Server
}

func New(addr string, paymentHandler *view.PaymentHandler) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/payments", paymentHandler.CreatePayment)

	return &Server{
		httpServer: &http.Server{Addr: addr, Handler: mux},
	}
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}
