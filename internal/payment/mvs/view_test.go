package mvs

import "testing"

func TestPaymentPostView_ToModel_And_FromModel(t *testing.T) {
	v := NewPaymentPostView()

	in := []byte(`{"amount":250,"currency":"EUR","order_id":"ord_9"}`)
	model, err := v.ToModel(in)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if model.Amount != 250 || model.Currency != "EUR" || model.OrderID != "ord_9" {
		t.Fatalf("unexpected model: %+v", model)
	}

	outBytes, err := v.FromModel(PaymentPostResponse{Status: "ok", PaymentID: "p_1", Message: "m"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if string(outBytes) == "" {
		t.Fatalf("expected non-empty json")
	}
}
