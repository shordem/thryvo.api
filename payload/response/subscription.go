package response

type PaymentInitialization struct {
	Reference        string `json:"reference"`
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"access_code"`
}

type SubscriptionPlan struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Duration    int     `json:"duration"`
}

type UserSubscription struct {
	ID            string           `json:"id"`
	Plan          SubscriptionPlan `json:"plan"`
	Status        string           `json:"status"`
	StartDate     int64            `json:"start_date"`
	EndDate       int64            `json:"end_date"`
	IsActive      bool             `json:"is_active"`
	DaysRemaining int              `json:"days_remaining"`
}

type Transaction struct {
	ID               string  `json:"id"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	PaymentReference string  `json:"payment_reference"`
	PaymentMethod    string  `json:"payment_method"`
	Status           string  `json:"status"`
	CreatedAt        int64   `json:"created_at"`
}

type SubscriptionStatus struct {
	HasSubscription bool   `json:"has_subscription"`
	Status          string `json:"status"`
	DaysRemaining   int    `json:"days_remaining"`
	EndDate         int64  `json:"end_date"`
}
