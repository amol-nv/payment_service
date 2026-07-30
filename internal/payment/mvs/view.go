package mvs

import "encoding/json"

// PaymentPostView defines the View mapping between HTTP/JSON representations and MVS Models.
// In this repo, the View is responsible for JSON marshal/unmarshal and field mapping.
//
// Note: The handler can use these helpers to map request/response payloads.

type PaymentPostView struct{}

func NewPaymentPostView() *PaymentPostView {
	return &PaymentPostView{}
}

// ToModel maps an incoming JSON payload (HTTP request body) into the MVS Model.
func (v *PaymentPostView) ToModel(payload []byte) (PaymentPostRequest, error) {
	var req PaymentPostRequest
	if err := json.Unmarshal(payload, &req); err != nil {
		return PaymentPostRequest{}, err
	}
	return req, nil
}

// FromModel maps the MVS Model into a JSON payload (HTTP response body).
func (v *PaymentPostView) FromModel(model PaymentPostResponse) ([]byte, error) {
	return json.Marshal(model)
}
