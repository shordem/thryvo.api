package payment

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	PaystackBaseURL = "https://api.paystack.co"
)

type PaystackGateway struct {
	secretKey string
	client    *http.Client
}

func NewPaystackGateway(secretKey string) *PaystackGateway {
	return &PaystackGateway{
		secretKey: secretKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *PaystackGateway) GetName() string {
	return "paystack"
}

// Initialize initializes a payment transaction
func (p *PaystackGateway) Initialize(ctx context.Context, req *InitializePaymentRequest) (*InitializePaymentResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", PaystackBaseURL)

	// Convert amount to kobo (Paystack uses smallest currency unit)
	amountInKobo := int64(req.Amount * 100)

	payload := map[string]interface{}{
		"email":        req.Email,
		"amount":       amountInKobo,
		"reference":    req.Reference,
		"callback_url": req.CallbackURL,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			AuthorizationURL string `json:"authorization_url"`
			AccessCode       string `json:"access_code"`
			Reference        string `json:"reference"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, errors.New("payment initialization failed: " + result.Message)
	}

	return &InitializePaymentResponse{
		Reference:        result.Data.Reference,
		AuthorizationURL: result.Data.AuthorizationURL,
		AccessCode:       result.Data.AccessCode,
	}, nil
}

// VerifyPayment verifies a payment transaction
func (p *PaystackGateway) VerifyPayment(ctx context.Context, reference string) (*VerificationResponse, error) {
	url := fmt.Sprintf("%s/transaction/verify/%s", PaystackBaseURL, reference)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.secretKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Reference string `json:"reference"`
			Amount    int64  `json:"amount"`
			Currency  string `json:"currency"`
			Status    string `json:"status"`
			PaidAt    string `json:"paid_at"`
			Channel   string `json:"channel"`
			Customer  struct {
				Email string `json:"email"`
			} `json:"customer"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, errors.New("payment verification failed: " + result.Message)
	}

	// Convert amount back from kobo to naira
	amount := float64(result.Data.Amount) / 100

	return &VerificationResponse{
		Reference: result.Data.Reference,
		Amount:    amount,
		Currency:  result.Data.Currency,
		Status:    result.Data.Status,
		PaidAt:    result.Data.PaidAt,
		Channel:   result.Data.Channel,
		Customer: CustomerData{
			Email: result.Data.Customer.Email,
		},
	}, nil
}

// ValidateWebhook validates Paystack webhook signature
func (p *PaystackGateway) ValidateWebhook(signature string, payload []byte) bool {
	mac := hmac.New(sha512.New, []byte(p.secretKey))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
