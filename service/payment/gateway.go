package payment

import (
	"context"
	"errors"
)

// PaymentGateway defines the interface for payment gateways
type PaymentGateway interface {
	Initialize(ctx context.Context, req *InitializePaymentRequest) (*InitializePaymentResponse, error)
	VerifyPayment(ctx context.Context, reference string) (*VerificationResponse, error)
	ValidateWebhook(signature string, payload []byte) bool
	GetName() string
}

// InitializePaymentRequest contains payment initialization data
type InitializePaymentRequest struct {
	Email       string
	Amount      float64
	Currency    string
	Reference   string
	CallbackURL string
}

// InitializePaymentResponse contains payment initialization response
type InitializePaymentResponse struct {
	Reference        string
	AuthorizationURL string
	AccessCode       string
}

// VerificationResponse contains payment verification data
type VerificationResponse struct {
	Reference string
	Amount    float64
	Currency  string
	Status    string
	PaidAt    string
	Channel   string
	Customer  CustomerData
}

type CustomerData struct {
	Email string
}

// PaymentGatewayFactory creates payment gateway instances
type PaymentGatewayFactory struct {
	gateways map[string]PaymentGateway
}

func NewPaymentGatewayFactory() *PaymentGatewayFactory {
	return &PaymentGatewayFactory{
		gateways: make(map[string]PaymentGateway),
	}
}

func (f *PaymentGatewayFactory) Register(name string, gateway PaymentGateway) {
	f.gateways[name] = gateway
}

func (f *PaymentGatewayFactory) Get(name string) (PaymentGateway, error) {
	gateway, exists := f.gateways[name]
	if !exists {
		return nil, errors.New("payment gateway not found: " + name)
	}
	return gateway, nil
}

func (f *PaymentGatewayFactory) GetAvailable() []string {
	names := make([]string, 0, len(f.gateways))
	for name := range f.gateways {
		names = append(names, name)
	}
	return names
}
