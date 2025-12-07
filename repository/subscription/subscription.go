package subscription_repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/dto"
	"github.com/shordem/api.thryvo/lib/database"
	"github.com/shordem/api.thryvo/model"
)

type SubscriptionRepositoryInterface interface {
	// Subscription Plans
	GetAllPlans(ctx context.Context) ([]dto.SubscriptionPlan, error)
	GetPlanByID(ctx context.Context, id uuid.UUID) (*dto.SubscriptionPlan, error)
	CreatePlan(ctx context.Context, plan *model.SubscriptionPlan) error
	UpdatePlan(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeletePlan(ctx context.Context, id uuid.UUID) error

	// User Subscriptions
	CreateSubscription(ctx context.Context, subscription *model.UserSubscription) error
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*dto.UserSubscription, error)
	GetUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]dto.UserSubscription, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	ExpireOldSubscriptions(ctx context.Context) error

	// Transactions
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]dto.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id uuid.UUID, status string) error
}

type subscriptionRepository struct {
	database database.DatabaseInterface
}

func NewSubscriptionRepository(database database.DatabaseInterface) SubscriptionRepositoryInterface {
	return &subscriptionRepository{database: database}
}

// Subscription Plans
func (r *subscriptionRepository) GetAllPlans(ctx context.Context) ([]dto.SubscriptionPlan, error) {
	var plans []model.SubscriptionPlan
	err := r.database.Connection().WithContext(ctx).
		Where("is_active = ? AND deleted_at IS NULL", true).
		Find(&plans).Error

	if err != nil {
		return nil, err
	}

	result := make([]dto.SubscriptionPlan, len(plans))
	for i, plan := range plans {
		result[i] = r.planToDTO(&plan)
	}

	return result, nil
}

func (r *subscriptionRepository) GetPlanByID(ctx context.Context, id uuid.UUID) (*dto.SubscriptionPlan, error) {
	var plan model.SubscriptionPlan
	err := r.database.Connection().WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&plan).Error

	if err != nil {
		return nil, err
	}

	result := r.planToDTO(&plan)
	return &result, nil
}

func (r *subscriptionRepository) CreatePlan(ctx context.Context, plan *model.SubscriptionPlan) error {
	plan.Prepare()
	return r.database.Connection().WithContext(ctx).Create(plan).Error
}

func (r *subscriptionRepository) UpdatePlan(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.database.Connection().WithContext(ctx).
		Model(&model.SubscriptionPlan{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *subscriptionRepository) DeletePlan(ctx context.Context, id uuid.UUID) error {
	now := time.Now().Unix()
	return r.database.Connection().WithContext(ctx).
		Model(&model.SubscriptionPlan{}).
		Where("id = ?", id).
		Update("deleted_at", now).Error
}

// User Subscriptions
func (r *subscriptionRepository) CreateSubscription(ctx context.Context, subscription *model.UserSubscription) error {
	return r.database.Connection().WithContext(ctx).Create(subscription).Error
}

func (r *subscriptionRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*dto.UserSubscription, error) {
	var subscription model.UserSubscription
	err := r.database.Connection().WithContext(ctx).
		Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, "active").
		First(&subscription).Error

	if err != nil {
		return nil, err
	}

	result := r.subscriptionToDTO(&subscription)
	return &result, nil
}

func (r *subscriptionRepository) GetUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]dto.UserSubscription, error) {
	var subscriptions []model.UserSubscription
	err := r.database.Connection().WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&subscriptions).Error

	if err != nil {
		return nil, err
	}

	result := make([]dto.UserSubscription, len(subscriptions))
	for i, sub := range subscriptions {
		result[i] = r.subscriptionToDTO(&sub)
	}

	return result, nil
}

func (r *subscriptionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.database.Connection().WithContext(ctx).
		Model(&model.UserSubscription{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *subscriptionRepository) ExpireOldSubscriptions(ctx context.Context) error {
	now := time.Now()
	return r.database.Connection().WithContext(ctx).
		Model(&model.UserSubscription{}).
		Where("status = ? AND end_date < ?", "active", now).
		Update("status", "expired").Error
}

// Transactions
func (r *subscriptionRepository) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	return r.database.Connection().WithContext(ctx).Create(transaction).Error
}

func (r *subscriptionRepository) GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]dto.Transaction, error) {
	var transactions []model.Transaction
	err := r.database.Connection().WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	result := make([]dto.Transaction, len(transactions))
	for i, txn := range transactions {
		result[i] = r.transactionToDTO(&txn)
	}

	return result, nil
}

func (r *subscriptionRepository) UpdateTransactionStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.database.Connection().WithContext(ctx).
		Model(&model.Transaction{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// DTO Converters
func (r *subscriptionRepository) planToDTO(plan *model.SubscriptionPlan) dto.SubscriptionPlan {
	return dto.SubscriptionPlan{
		ID:          plan.ID,
		Name:        plan.Name,
		Description: plan.Description,
		Price:       plan.Price,
		Currency:    plan.Currency,
		Duration:    plan.Duration,
		IsActive:    plan.IsActive,
		CreatedAt:   plan.CreatedAt,
		UpdatedAt:   plan.UpdatedAt,
	}
}

func (r *subscriptionRepository) subscriptionToDTO(subscription *model.UserSubscription) dto.UserSubscription {
	return dto.UserSubscription{
		ID:        subscription.ID,
		UserID:    subscription.UserID,
		PlanID:    subscription.PlanID,
		Status:    subscription.Status,
		StartDate: subscription.StartDate,
		EndDate:   subscription.EndDate,
		CreatedAt: subscription.CreatedAt,
		UpdatedAt: subscription.UpdatedAt,
	}
}

func (r *subscriptionRepository) transactionToDTO(transaction *model.Transaction) dto.Transaction {
	return dto.Transaction{
		ID:               transaction.ID,
		UserID:           transaction.UserID,
		SubscriptionID:   transaction.SubscriptionID,
		Amount:           transaction.Amount,
		Currency:         transaction.Currency,
		PaymentReference: transaction.PaymentReference,
		PaymentMethod:    transaction.PaymentMethod,
		Status:           transaction.Status,
		CreatedAt:        transaction.CreatedAt,
		UpdatedAt:        transaction.UpdatedAt,
	}
}
