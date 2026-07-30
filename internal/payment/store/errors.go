package store

import "fmt"

type HTTPError struct {
	StatusCode int
	Body       []byte
	Message    string
}

func (e *HTTPError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("payment store http error: status=%d message=%s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("payment store http error: status=%d", e.StatusCode)
}
