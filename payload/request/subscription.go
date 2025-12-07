package request

type InitializePayment struct {
	PlanID  string `json:"plan_id" validate:"required"`
	Gateway string `json:"gateway" validate:"required"` // paystack, flutterwave, etc
}

type VerifyPayment struct {
	Reference string `json:"reference" validate:"required"`
	Gateway   string `json:"gateway" validate:"required"`
}

type CreateSubscription struct {
	PlanID           string  `json:"plan_id" validate:"required"`
	PaymentReference string  `json:"payment_reference" validate:"required"`
	PaymentMethod    string  `json:"payment_method"`
	Amount           float64 `json:"amount" validate:"required,gt=0"`
	Currency         string  `json:"currency" validate:"required"`
}
