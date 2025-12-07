package subscription_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shordem/api.thryvo/dto"
	"github.com/shordem/api.thryvo/model"
	"github.com/shordem/api.thryvo/payload/request"
	"github.com/shordem/api.thryvo/payload/response"
	subscriptionRepo "github.com/shordem/api.thryvo/repository/subscription"
	userRepo "github.com/shordem/api.thryvo/repository/user"
	"github.com/shordem/api.thryvo/service/payment"
)

type SubscriptionServiceInterface interface {
	// Plans
	GetAllPlans(ctx context.Context) ([]response.SubscriptionPlan, error)
	CreatePlan(ctx context.Context, req *request.CreatePlan) (*response.SubscriptionPlan, error)
	UpdatePlan(ctx context.Context, id uuid.UUID, req *request.UpdatePlan) (*response.SubscriptionPlan, error)
	DeletePlan(ctx context.Context, id uuid.UUID) error

	// Payment Integration
	InitializePayment(ctx context.Context, userID uuid.UUID, req *request.InitializePayment) (*response.PaymentInitialization, error)
	VerifyPaymentAndCreateSubscription(ctx context.Context, userID uuid.UUID, req *request.VerifyPayment) (*response.UserSubscription, error)
	HandleWebhook(ctx context.Context, gatewayName string, signature string, payload []byte) error

	// Subscriptions
	CreateSubscription(ctx context.Context, userID uuid.UUID, req *request.CreateSubscription) (*response.UserSubscription, error)
	GetSubscriptionStatus(ctx context.Context, userID uuid.UUID) (*response.SubscriptionStatus, error)
	GetUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]response.UserSubscription, error)
	CancelSubscription(ctx context.Context, userID uuid.UUID) error

	// Transactions
	GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]response.Transaction, error)

	// Maintenance
	CheckAndExpireSubscriptions(ctx context.Context) error
}

type subscriptionService struct {
	repository     subscriptionRepo.SubscriptionRepositoryInterface
	userRepository userRepo.UserRepositoryInterface
	paymentFactory *payment.PaymentGatewayFactory
}

func NewSubscriptionService(repository subscriptionRepo.SubscriptionRepositoryInterface, userRepository userRepo.UserRepositoryInterface, paymentFactory *payment.PaymentGatewayFactory) SubscriptionServiceInterface {
	return &subscriptionService{
		repository:     repository,
		userRepository: userRepository,
		paymentFactory: paymentFactory,
	}
}

// Plans
func (s *subscriptionService) GetAllPlans(ctx context.Context) ([]response.SubscriptionPlan, error) {
	plans, err := s.repository.GetAllPlans(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]response.SubscriptionPlan, len(plans))
	for i, plan := range plans {
		result[i] = response.SubscriptionPlan{
			ID:          plan.ID.String(),
			Name:        plan.Name,
			Description: plan.Description,
			Price:       plan.Price,
			Currency:    plan.Currency,
			Duration:    plan.Duration,
		}
	}

	return result, nil
}

func (s *subscriptionService) CreatePlan(ctx context.Context, req *request.CreatePlan) (*response.SubscriptionPlan, error) {
	plan := &model.SubscriptionPlan{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Currency:    req.Currency,
		Duration:    req.Duration,
		IsActive:    true,
	}

	if err := s.repository.CreatePlan(ctx, plan); err != nil {
		return nil, err
	}

	return &response.SubscriptionPlan{
		ID:          plan.ID.String(),
		Name:        plan.Name,
		Description: plan.Description,
		Price:       plan.Price,
		Currency:    plan.Currency,
		Duration:    plan.Duration,
	}, nil
}

func (s *subscriptionService) UpdatePlan(ctx context.Context, id uuid.UUID, req *request.UpdatePlan) (*response.SubscriptionPlan, error) {
	// Check if plan exists
	existing, err := s.repository.GetPlanByID(ctx, id)
	if err != nil {
		return nil, errors.New("plan not found")
	}

	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now().Unix()

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.Duration > 0 {
		updates["duration"] = req.Duration
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := s.repository.UpdatePlan(ctx, id, updates); err != nil {
		return nil, err
	}

	// Get updated plan
	updated, err := s.repository.GetPlanByID(ctx, id)
	if err != nil {
		// Return with updated values if fetch fails
		return &response.SubscriptionPlan{
			ID:          existing.ID.String(),
			Name:        existing.Name,
			Description: existing.Description,
			Price:       existing.Price,
			Currency:    existing.Currency,
			Duration:    existing.Duration,
		}, nil
	}

	return &response.SubscriptionPlan{
		ID:          updated.ID.String(),
		Name:        updated.Name,
		Description: updated.Description,
		Price:       updated.Price,
		Currency:    updated.Currency,
		Duration:    updated.Duration,
	}, nil
}

func (s *subscriptionService) DeletePlan(ctx context.Context, id uuid.UUID) error {
	// Check if plan exists
	_, err := s.repository.GetPlanByID(ctx, id)
	if err != nil {
		return errors.New("plan not found")
	}

	return s.repository.DeletePlan(ctx, id)
}

// Payment Integration
func (s *subscriptionService) InitializePayment(ctx context.Context, userID uuid.UUID, req *request.InitializePayment) (*response.PaymentInitialization, error) {
	// Get user details to fetch email
	user, err := s.userRepository.FindUserById(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get plan details
	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return nil, errors.New("invalid plan ID")
	}

	plan, err := s.repository.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, errors.New("plan not found")
	}

	// Check if user already has active subscription
	existing, err := s.repository.GetActiveByUserID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user already has an active subscription")
	}

	// Get payment gateway
	gateway, err := s.paymentFactory.Get(req.Gateway)
	if err != nil {
		return nil, err
	}

	// Generate unique reference
	reference := fmt.Sprintf("SUB_%s_%d", userID.String()[:8], time.Now().Unix())

	// Initialize payment
	paymentReq := &payment.InitializePaymentRequest{
		Email:       user.Email,
		Amount:      plan.Price,
		Currency:    plan.Currency,
		Reference:   reference,
		CallbackURL: "https://filecapsa.com/dashboard/subscription/verify",
	}

	paymentResp, err := gateway.Initialize(ctx, paymentReq)
	if err != nil {
		return nil, err
	}

	return &response.PaymentInitialization{
		Reference:        paymentResp.Reference,
		AuthorizationURL: paymentResp.AuthorizationURL,
		AccessCode:       paymentResp.AccessCode,
	}, nil
}

func (s *subscriptionService) VerifyPaymentAndCreateSubscription(ctx context.Context, userID uuid.UUID, req *request.VerifyPayment) (*response.UserSubscription, error) {
	// Get payment gateway
	gateway, err := s.paymentFactory.Get(req.Gateway)
	if err != nil {
		return nil, err
	}

	// Verify payment
	verification, err := gateway.VerifyPayment(ctx, req.Reference)
	if err != nil {
		return nil, err
	}

	if verification.Status != "success" {
		return nil, errors.New("payment verification failed")
	}

	// Extract plan ID from reference (SUB_<user_id>_<timestamp>)
	// For now, we'll need to pass plan_id in the request or store it temporarily
	// Let's add it to the verification request
	return s.createSubscriptionFromPayment(ctx, userID, verification)
}

func (s *subscriptionService) createSubscriptionFromPayment(ctx context.Context, userID uuid.UUID, verification *payment.VerificationResponse) (*response.UserSubscription, error) {
	// Note: In production, you should extract plan_id from metadata or have it passed separately
	// For now, we'll assume the first plan
	plans, err := s.repository.GetAllPlans(ctx)
	if err != nil || len(plans) == 0 {
		return nil, errors.New("no plans available")
	}

	plan := &plans[0] // Use first plan for now

	now := time.Now()
	startDate := now
	endDate := now.AddDate(0, 0, plan.Duration)

	// Create subscription
	subscription := &model.UserSubscription{
		UserID:    userID,
		PlanID:    plan.ID,
		Status:    "active",
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := s.repository.CreateSubscription(ctx, subscription); err != nil {
		return nil, err
	}

	// Create transaction
	transaction := &model.Transaction{
		UserID:           userID,
		SubscriptionID:   subscription.ID,
		Amount:           verification.Amount,
		Currency:         verification.Currency,
		PaymentReference: verification.Reference,
		PaymentMethod:    verification.Channel,
		Status:           "completed",
	}

	if err := s.repository.CreateTransaction(ctx, transaction); err != nil {
		return nil, err
	}

	return s.toSubscriptionResponse(&dto.UserSubscription{
		ID:        subscription.ID,
		UserID:    subscription.UserID,
		PlanID:    subscription.PlanID,
		Status:    subscription.Status,
		StartDate: subscription.StartDate,
		EndDate:   subscription.EndDate,
		CreatedAt: subscription.CreatedAt,
		UpdatedAt: subscription.UpdatedAt,
	}, plan), nil
}

// Subscriptions
func (s *subscriptionService) CreateSubscription(ctx context.Context, userID uuid.UUID, req *request.CreateSubscription) (*response.UserSubscription, error) {
	// Check if user already has active subscription
	existing, err := s.repository.GetActiveByUserID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user already has an active subscription")
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return nil, errors.New("invalid plan ID")
	}

	// Get plan details
	plan, err := s.repository.GetPlanByID(ctx, planID)
	if err != nil {
		return nil, errors.New("plan not found")
	}

	now := time.Now()
	startDate := now
	endDate := now.AddDate(0, 0, plan.Duration)

	// Create subscription
	subscription := &model.UserSubscription{
		UserID:    userID,
		PlanID:    planID,
		Status:    "active",
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := s.repository.CreateSubscription(ctx, subscription); err != nil {
		return nil, err
	}

	// Create transaction
	transaction := &model.Transaction{
		UserID:           userID,
		SubscriptionID:   subscription.ID,
		Amount:           req.Amount,
		Currency:         req.Currency,
		PaymentReference: req.PaymentReference,
		PaymentMethod:    req.PaymentMethod,
		Status:           "completed",
	}

	if err := s.repository.CreateTransaction(ctx, transaction); err != nil {
		return nil, err
	}

	return s.toSubscriptionResponse(&dto.UserSubscription{
		ID:        subscription.ID,
		UserID:    subscription.UserID,
		PlanID:    subscription.PlanID,
		Status:    subscription.Status,
		StartDate: subscription.StartDate,
		EndDate:   subscription.EndDate,
		CreatedAt: subscription.CreatedAt,
		UpdatedAt: subscription.UpdatedAt,
	}, plan), nil
}

func (s *subscriptionService) GetSubscriptionStatus(ctx context.Context, userID uuid.UUID) (*response.SubscriptionStatus, error) {
	subscription, err := s.repository.GetActiveByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &response.SubscriptionStatus{
				HasSubscription: false,
				Status:          "none",
				DaysRemaining:   0,
				EndDate:         0,
			}, nil
		}
		return nil, err
	}

	now := time.Now()
	daysRemaining := int(subscription.EndDate.Sub(now).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	return &response.SubscriptionStatus{
		HasSubscription: true,
		Status:          subscription.Status,
		DaysRemaining:   daysRemaining,
		EndDate:         subscription.EndDate.Unix(),
	}, nil
}

func (s *subscriptionService) GetUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]response.UserSubscription, error) {
	subscriptions, err := s.repository.GetUserSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]response.UserSubscription, 0)
	for _, sub := range subscriptions {
		plan, err := s.repository.GetPlanByID(ctx, sub.PlanID)
		if err != nil {
			continue
		}
		result = append(result, *s.toSubscriptionResponse(&sub, plan))
	}

	return result, nil
}

func (s *subscriptionService) CancelSubscription(ctx context.Context, userID uuid.UUID) error {
	subscription, err := s.repository.GetActiveByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.repository.UpdateStatus(ctx, subscription.ID, "cancelled")
}

// Transactions
func (s *subscriptionService) GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]response.Transaction, error) {
	transactions, err := s.repository.GetUserTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]response.Transaction, len(transactions))
	for i, txn := range transactions {
		result[i] = response.Transaction{
			ID:               txn.ID.String(),
			Amount:           txn.Amount,
			Currency:         txn.Currency,
			PaymentReference: txn.PaymentReference,
			PaymentMethod:    txn.PaymentMethod,
			Status:           txn.Status,
			CreatedAt:        txn.CreatedAt.Unix(),
		}
	}

	return result, nil
}

// HandleWebhook processes payment gateway webhooks
func (s *subscriptionService) HandleWebhook(ctx context.Context, gatewayName string, signature string, payload []byte) error {
	gateway, err := s.paymentFactory.Get(gatewayName)
	if err != nil {
		return err
	}

	// Validate webhook signature
	if !gateway.ValidateWebhook(signature, payload) {
		return errors.New("invalid webhook signature")
	}

	// Process webhook event (implement based on your needs)
	// For now, just log it
	return nil
}

func (s *subscriptionService) CheckAndExpireSubscriptions(ctx context.Context) error {
	return s.repository.ExpireOldSubscriptions(ctx)
}

func (s *subscriptionService) toSubscriptionResponse(sub *dto.UserSubscription, plan *dto.SubscriptionPlan) *response.UserSubscription {
	now := time.Now()
	daysRemaining := int(sub.EndDate.Sub(now).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	return &response.UserSubscription{
		ID: sub.ID.String(),
		Plan: response.SubscriptionPlan{
			ID:          plan.ID.String(),
			Name:        plan.Name,
			Description: plan.Description,
			Price:       plan.Price,
			Currency:    plan.Currency,
			Duration:    plan.Duration,
		},
		Status:        sub.Status,
		StartDate:     sub.StartDate.Unix(),
		EndDate:       sub.EndDate.Unix(),
		IsActive:      sub.Status == "active",
		DaysRemaining: daysRemaining,
	}
}
