package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/handler"
	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/lib/database"
	subscriptionService "github.com/shordem/api.thryvo/service/subscription"
	user_service "github.com/shordem/api.thryvo/service/user"
)

// subscriptionAdapter adapts the subscription service to the SubscriptionChecker interface
type subscriptionAdapter struct {
	subService subscriptionService.SubscriptionServiceInterface
}

func (s *subscriptionAdapter) GetSubscriptionStatus(ctx context.Context, userID uuid.UUID) (*user_service.SubscriptionStatusResponse, error) {
	status, err := s.subService.GetSubscriptionStatus(ctx, userID)
	if err != nil {
		return &user_service.SubscriptionStatusResponse{IsActive: false}, err
	}
	isActive := status.HasSubscription && status.Status == "active"
	return &user_service.SubscriptionStatusResponse{IsActive: isActive}, nil
}

func InitializeRouter(router *fiber.App, dbConn database.DatabaseInterface, env constants.Env) {

	main := router.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	main.Get("/monitor", monitor.New(monitor.Config{Title: "Thryvo API Monitor"}))

	// Initialize subscription router first to get the service
	subscriptionService := InitializeSubscriptionRouter(main, dbConn, env)

	// Create adapter for subscription service
	subAdapter := &subscriptionAdapter{subService: subscriptionService}

	InitializeUserRouter(main, dbConn, env, subAdapter)
	InitializeCoreRouter(main, dbConn, env)

	router.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	router.Get("/logs/:key", func(c *fiber.Ctx) error {
		return handler.GetLogs(c, dbConn)
	})
	router.Get("/", handler.Index)
	router.Get("*", handler.NotFound)

}
