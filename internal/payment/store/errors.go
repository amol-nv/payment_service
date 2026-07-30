package store

import "errors"

var (
	ErrBadRequest      = errors.New("bad request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrNotFound        = errors.New("not found")
	ErrConflict        = errors.New("conflict")
	ErrUnprocessable   = errors.New("unprocessable entity")
	ErrUpstreamFailure = errors.New("upstream failure")
)
