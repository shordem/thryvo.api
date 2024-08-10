package user_handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/handler"
	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/lib/helper"
	"github.com/shordem/api.thryvo/payload/response"
	userResponse "github.com/shordem/api.thryvo/payload/response/user"
	userService "github.com/shordem/api.thryvo/service/user"
)

type userHandler struct {
	userService   userService.UserServiceInterface
	authconstants helper.AuthInterface
}

type UserHandlerInterface interface {
	UserDetails(c *fiber.Ctx) error
	FindAllUsers(c *fiber.Ctx) error
}

func NewUserHandler(userService userService.UserServiceInterface) UserHandlerInterface {
	return &userHandler{
		userService:   userService,
		authconstants: helper.NewAuth(),
	}
}

func (u *userHandler) UserDetails(c *fiber.Ctx) error {
	var resp userResponse.UserResponse

	userId := c.Locals("userId").(uuid.UUID)

	user, err := u.userService.FindUserById(userId.String())

	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Customer details retrieved successfully"
	resp.Data.FirstName = user.FirstName
	resp.Data.LastName = user.LastName
	resp.Data.Email = user.Email

	return c.Status(http.StatusOK).JSON(resp)
}

// FindAllUsers is a method that returns a list of all users
// Role: Admin
func (u *userHandler) FindAllUsers(c *fiber.Ctx) error {
	var resp response.Response
	pageable := handler.GeneratePageable(c)

	users, pagination, err := u.userService.FindAllUsers(pageable)

	if err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = "Users retrieved successfully"
	resp.Data = map[string]interface{}{"results": users, "pagination": pagination}

	return c.Status(http.StatusOK).JSON(resp)
}
