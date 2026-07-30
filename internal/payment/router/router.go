package router

import (
	"net/http"

	"amol-nv/payment_service/internal/payment/service"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(paymentSvc *service.PaymentService) *Router {
	m := http.NewServeMux()
	m.HandleFunc("/api/payments", paymentSvc.HandlePostPayment)
	return &Router{mux: m}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
