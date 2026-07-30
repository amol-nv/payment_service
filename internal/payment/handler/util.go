package handler

import (
	"errors"
	"io"
	"net/http"
)

func readAllLimit(r *http.Request, limit int64) ([]byte, error) {
	if r.Body == nil {
		return nil, errors.New("empty body")
	}
	defer r.Body.Close()

	lr := &io.LimitedReader{R: r.Body, N: limit + 1}
	b, err := io.ReadAll(lr)
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > limit {
		return nil, errors.New("body too large")
	}
	return b, nil
}
