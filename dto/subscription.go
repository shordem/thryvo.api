package dto

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionPlan struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	Duration    int       `json:"duration"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserSubscription struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	PlanID    uuid.UUID `json:"plan_id"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	SubscriptionID   uuid.UUID `json:"subscription_id"`
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	PaymentReference string    `json:"payment_reference"`
	PaymentMethod    string    `json:"payment_method"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
