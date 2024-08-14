package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/shordem/api.thryvo/lib/database"
	user_repository "github.com/shordem/api.thryvo/repository/user"
)

func RequireAPIKey(db database.DatabaseInterface) fiber.Handler {
	keyRepo := user_repository.NewKeyRepository(db)

	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-KEY")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "API key is required"})
		}

		key, err := keyRepo.FindUserIDByKey(apiKey)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid API key"})
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
		}

		c.Locals("userId", key.UserID)

		return c.Next()
	}
}
