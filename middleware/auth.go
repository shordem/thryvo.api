package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/shordem/api.thryvo/lib/helper"
)

func Protected() fiber.Handler {
	authHelper := helper.NewAuth()

	return func(c *fiber.Ctx) (err error) {
		token := authHelper.ExtractBearerToken(c.Request())
		userId, err := authHelper.ExtractUserID(token, "access")

		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
		}

		c.Locals("userId", userId)

		return c.Next()
	}
}
