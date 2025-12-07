package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/lib/database"
)

type User struct {
	database.BaseModel

	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"is_email_verified"`
	Password        string `json:"password"`
	Role            string `json:"role"`
}

type VerificationCode struct {
	database.BaseModel

	UserID  uuid.UUID `json:"user_id"`
	Code    string    `json:"code"`
	Purpose string    `json:"purpose"`
}

type Key struct {
	database.BaseModel

	UserID uuid.UUID `json:"user_id"`
	Key    string    `json:"key"`
}

type SubscriptionPlan struct {
	database.BaseModel

	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Duration    int     `json:"duration"` // in days
	IsActive    bool    `json:"is_active"`
}

type UserSubscription struct {
	database.BaseModel

	UserID    uuid.UUID `json:"user_id"`
	PlanID    uuid.UUID `json:"plan_id"`
	Status    string    `json:"status"` // active, expired, cancelled
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type Transaction struct {
	database.BaseModel

	UserID           uuid.UUID `json:"user_id"`
	SubscriptionID   uuid.UUID `json:"subscription_id"`
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	PaymentReference string    `json:"payment_reference"`
	PaymentMethod    string    `json:"payment_method"`
	Status           string    `json:"status"` // pending, completed, failed
}
