package service

import "errors"

var (
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrInvalidCurrency    = errors.New("invalid currency")
	ErrInvalidMethod      = errors.New("invalid method")
	ErrInvalidOrderID     = errors.New("invalid order id")
	ErrPaymentBadRequest = errors.New("payment bad request")
	ErrPaymentUnauthorized = errors.New("payment unauthorized")
	ErrPaymentForbidden    = errors.New("payment forbidden")
	ErrPaymentNotFound     = errors.New("payment not found")
	ErrPaymentConflict     = errors.New("payment conflict")
	ErrPaymentUnprocessable = errors.New("payment unprocessable")
	ErrPaymentUpstreamFailure = errors.New("payment upstream failure")
)
