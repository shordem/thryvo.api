package router

import (
	"github.com/gofiber/fiber/v2"

	subscriptionHandler "github.com/shordem/api.thryvo/handler/subscription"
	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/lib/database"
	"github.com/shordem/api.thryvo/middleware"
	subscriptionRepo "github.com/shordem/api.thryvo/repository/subscription"
	user_repository "github.com/shordem/api.thryvo/repository/user"
	"github.com/shordem/api.thryvo/service/payment"
	subscriptionService "github.com/shordem/api.thryvo/service/subscription"
	user_service "github.com/shordem/api.thryvo/service/user"
)

func InitializeSubscriptionRouter(router fiber.Router, db database.DatabaseInterface, env constants.Env) subscriptionService.SubscriptionServiceInterface {
	// Payment Gateway Factory
	paymentFactory := payment.NewPaymentGatewayFactory()

	// Register Paystack gateway
	paystackGateway := payment.NewPaystackGateway(env.PAYSTACK_SECRET_KEY)
	paymentFactory.Register("paystack", paystackGateway)

	// Can add more gateways here
	// flutterwaveGateway := payment.NewFlutterwaveGateway(env.FLUTTERWAVE_SECRET_KEY)
	// paymentFactory.Register("flutterwave", flutterwaveGateway)

	// Repository
	subscriptionRepository := subscriptionRepo.NewSubscriptionRepository(db)
	userRepository := user_repository.NewUserRepository(db)

	// Service
	subService := subscriptionService.NewSubscriptionService(subscriptionRepository, userRepository, paymentFactory)

	// Handler
	subHandler := subscriptionHandler.NewHandler(subService)

	// Middleware
	authMiddleware := middleware.Protected()
	userRepo := user_repository.NewUserRepository(db)
	roleMiddleware := middleware.NewRoleMiddleware(userRepo)

	// Public routes
	router.Get("/subscriptions/plans", subHandler.GetPlans)

	// Protected routes
	subscriptionRoute := router.Group("/subscriptions", authMiddleware)

	subscriptionRoute.Post("/initialize-payment", subHandler.InitializePayment)
	subscriptionRoute.Post("/verify-payment", subHandler.VerifyPayment)
	subscriptionRoute.Get("/status", subHandler.GetSubscriptionStatus)
	subscriptionRoute.Get("/", subHandler.GetUserSubscriptions)
	subscriptionRoute.Get("/transactions", subHandler.GetTransactions)
	subscriptionRoute.Delete("/", subHandler.CancelSubscription)

	// Admin routes for plan management
	adminPlanRoute := router.Group("/subscriptions/plans", authMiddleware, roleMiddleware.ValidateRole(user_service.UserRoleAdmin))
	adminPlanRoute.Post("/", subHandler.CreatePlan)
	adminPlanRoute.Put("/:id", subHandler.UpdatePlan)
	adminPlanRoute.Delete("/:id", subHandler.DeletePlan)

	// Webhook endpoint (no auth required)
	router.Post("/subscriptions/webhook/:gateway", subHandler.PaymentWebhook)

	return subService
}
