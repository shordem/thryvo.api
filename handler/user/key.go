package user_handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/payload/response"
	user_service "github.com/shordem/api.thryvo/service/user"
)

type keyHandler struct {
	keyService user_service.KeyServiceInterface
}

type KeyServiceInterface interface {
	GetKey() string
}

func NewKeyHandler(keyService user_service.KeyServiceInterface) *keyHandler {
	return &keyHandler{keyService: keyService}
}

func (h *keyHandler) GetKey(c *fiber.Ctx) error {
	var resp response.Response

	userId := c.Locals("userId").(uuid.UUID)

	key, err := h.keyService.GetKey(userId)

	if err != nil {
		resp.Status = 400
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = "Key retrieved successfully"
	resp.Data = map[string]interface{}{"result": key}

	return c.Status(http.StatusOK).JSON(resp)
}
