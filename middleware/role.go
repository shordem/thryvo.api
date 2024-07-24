package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	user_repository "github.com/shordem/api.thryvo/repository/user"
)

type roleMiddleware struct {
	userRepository user_repository.UserRepositoryInterface
}

type RoleMiddlewareInterface interface {
	ValidateRole(role string) fiber.Handler
}

func NewRoleMiddleware(userRepository user_repository.UserRepositoryInterface) RoleMiddlewareInterface {
	return roleMiddleware{
		userRepository: userRepository,
	}
}

func (rm roleMiddleware) ValidateRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from JWT
		userID := c.Locals("userId").(uuid.UUID)

		user, err := rm.userRepository.FindUserById(userID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		if user.Role != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Forbidden",
			})
		}

		return c.Next()
	}
}
