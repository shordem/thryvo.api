package subscription

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/payload/request"
	"github.com/shordem/api.thryvo/payload/response"
	subscriptionService "github.com/shordem/api.thryvo/service/subscription"
)

type Handler struct {
	service subscriptionService.SubscriptionServiceInterface
}

func NewHandler(service subscriptionService.SubscriptionServiceInterface) *Handler {
	return &Handler{service: service}
}

// GetPlans gets all available subscription plans
// GET /api/subscriptions/plans
func (h *Handler) GetPlans(c *fiber.Ctx) error {
	var resp response.Response

	plans, err := h.service.GetAllPlans(c.Context())
	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Failed to get plans"
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Plans retrieved"
	resp.Data = map[string]interface{}{"plans": plans}

	return c.Status(http.StatusOK).JSON(resp)
}

// CreatePlan creates a new subscription plan (Admin only)
// POST /api/subscriptions/plans
func (h *Handler) CreatePlan(c *fiber.Ctx) error {
	var resp response.Response

	var req request.CreatePlan
	if err := c.BodyParser(&req); err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Invalid request body"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	plan, err := h.service.CreatePlan(c.Context(), &req)
	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Failed to create plan: " + err.Error()
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Plan created successfully"
	resp.Data = map[string]interface{}{"plan": plan}

	return c.Status(http.StatusCreated).JSON(resp)
}

// UpdatePlan updates an existing subscription plan (Admin only)
// PUT /api/subscriptions/plans/:id
func (h *Handler) UpdatePlan(c *fiber.Ctx) error {
	var resp response.Response

	planID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Invalid plan ID"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	var req request.UpdatePlan
	if err := c.BodyParser(&req); err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Invalid request body"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	plan, err := h.service.UpdatePlan(c.Context(), planID, &req)
	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Failed to update plan: " + err.Error()
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Plan updated successfully"
	resp.Data = map[string]interface{}{"plan": plan}

	return c.Status(http.StatusOK).JSON(resp)
}

// DeletePlan soft deletes a subscription plan (Admin only)
// DELETE /api/subscriptions/plans/:id
func (h *Handler) DeletePlan(c *fiber.Ctx) error {
	var resp response.Response

	planID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Invalid plan ID"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if err := h.service.DeletePlan(c.Context(), planID); err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Failed to delete plan: " + err.Error()
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Plan deleted successfully"
	resp.Data = map[string]interface{}{}

	return c.Status(http.StatusOK).JSON(resp)
}

// InitializePayment initializes a payment for subscription
// POST /api/subscriptions/initialize-payment
func (h *Handler) InitializePayment(c *fiber.Ctx) error {
	var resp response.Response

	userID := c.Locals("userId").(uuid.UUID)

	var req request.InitializePayment
	if err := c.BodyParser(&req); err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Invalid request body"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	payment, err := h.service.InitializePayment(c.Context(), userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "active subscription") {
			resp.Status = constants.ClientRequestValidationError
			resp.Message = err.Error()
			return c.Status(http.StatusConflict).JSON(resp)
		}
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Failed to initialize payment: " + err.Error()
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Payment initialized"
	resp.Data = map[string]interface{}{"payment": payment}

	return c.Status(http.StatusOK).JSON(resp)
}

// VerifyPayment verifies a payment and creates subscription
// POST /api/subscriptions/verify-payment
func (h *Handler) VerifyPayment(c *fiber.Ctx) error {
	var resp response.Response

	userID := c.Locals("userId").(uuid.UUID)

	var req request.VerifyPayment
	if err := c.BodyParser(&req); err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Invalid request body"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	subscription, err := h.service.VerifyPaymentAndCreateSubscription(c.Context(), userID, &req)
	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Failed to verify payment: " + err.Error()
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Subscription created successfully"
	resp.Data = map[string]interface{}{"subscription": subscription}

	return c.Status(http.StatusCreated).JSON(resp)
}

// PaymentWebhook handles payment gateway webhooks
// POST /api/subscriptions/webhook
func (h *Handler) PaymentWebhook(c *fiber.Ctx) error {
	var resp response.Response

	signature := c.Get("X-Paystack-Signature")
	body := c.Body()

	if err := h.service.HandleWebhook(c.Context(), "paystack", signature, body); err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Webhook processing failed: " + err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	return c.SendStatus(http.StatusOK)
}

// GetSubscriptionStatus gets current subscription status for a user
// GET /api/subscriptions/status
func (h *Handler) GetSubscriptionStatus(c *fiber.Ctx) error {
	var resp response.Response

	userID := c.Locals("userId").(uuid.UUID)

	status, err := h.service.GetSubscriptionStatus(c.Context(), userID)
	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Failed to get subscription status"
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Subscription status retrieved"
	resp.Data = map[string]interface{}{"status": status}

	return c.Status(http.StatusOK).JSON(resp)
}

// GetUserSubscriptions gets all subscriptions for a user
// GET /api/subscriptions
func (h *Handler) GetUserSubscriptions(c *fiber.Ctx) error {
	var resp response.Response

	userID := c.Locals("userId").(uuid.UUID)

	subscriptions, err := h.service.GetUserSubscriptions(c.Context(), userID)
	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Failed to get subscriptions"
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Subscriptions retrieved"
	resp.Data = map[string]interface{}{"subscriptions": subscriptions}

	return c.Status(http.StatusOK).JSON(resp)
}

// CancelSubscription cancels a user's active subscription
// POST /api/subscriptions/cancel
func (h *Handler) CancelSubscription(c *fiber.Ctx) error {
	var resp response.Response

	userID := c.Locals("userId").(uuid.UUID)

	if err := h.service.CancelSubscription(c.Context(), userID); err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Failed to cancel subscription"
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Subscription cancelled successfully"
	resp.Data = map[string]interface{}{}

	return c.Status(http.StatusOK).JSON(resp)
}

// GetTransactions gets all transactions for a user
// GET /api/subscriptions/transactions
func (h *Handler) GetTransactions(c *fiber.Ctx) error {
	var resp response.Response

	userID := c.Locals("userId").(uuid.UUID)

	transactions, err := h.service.GetUserTransactions(c.Context(), userID)
	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = "Failed to get transactions"
		return c.Status(http.StatusInternalServerError).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Transactions retrieved"
	resp.Data = map[string]interface{}{"transactions": transactions}

	return c.Status(http.StatusOK).JSON(resp)
}
