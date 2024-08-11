package user_handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/shordem/api.thryvo/dto"
	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/payload/request"
	"github.com/shordem/api.thryvo/payload/response"
	userResponse "github.com/shordem/api.thryvo/payload/response/user"
	userService "github.com/shordem/api.thryvo/service/user"
	"github.com/shordem/api.thryvo/validator"
)

type authHandler struct {
	authService userService.AuthServiceInterface
	validator   validator.AuthValidator
}

type AuthHandlerInterface interface {
	CheckEmail(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	RefreshAccessToken(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	ResendEmailVerification(c *fiber.Ctx) error
	VerifyEmail(c *fiber.Ctx) error
	VerifyEmailAndCode(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
}

func NewAuthHandler(authService userService.AuthServiceInterface) AuthHandlerInterface {
	return &authHandler{authService: authService}
}

func (handler *authHandler) CheckEmail(c *fiber.Ctx) error {
	var resp response.Response

	emailRequest := new(request.EmailRequest)

	if err := c.BodyParser(emailRequest); err != nil {
		resp.Status = constants.ClientUnProcessableEntity
		resp.Message = http.StatusText(http.StatusUnprocessableEntity)
		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if vEs, err := handler.validator.EmailValidate(*emailRequest); err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = err.Error()
		resp.Data = vEs
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	status, err := handler.authService.CheckEmail(emailRequest.Email)

	if err != nil {
		resp.Status = status
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = status
	resp.Message = http.StatusText(http.StatusOK)

	return c.JSON(resp)
}

func (handler *authHandler) Login(c *fiber.Ctx) error {
	var resp userResponse.LoginResponse

	loginRequest := new(request.LoginRequest)

	if err := c.BodyParser(loginRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if _, err := handler.validator.LoginValidate(*loginRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	token, status, err := handler.authService.Login(loginRequest.Email, loginRequest.Password)

	if err != nil {
		resp.Status = status
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = status
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = token

	return c.JSON(resp)
}

func (handler *authHandler) Register(c *fiber.Ctx) error {
	var resp response.Response
	var authDto dto.AuthDTO

	registerRequest := new(request.RegisterRequest)

	if err := c.BodyParser(registerRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if vEs, err := handler.validator.RegisterValidate(*registerRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		resp.Data = vEs
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	authDto.FirstName = registerRequest.FirstName
	authDto.LastName = registerRequest.LastName
	authDto.Email = registerRequest.Email
	authDto.Password = registerRequest.Password

	if err := handler.authService.Register(authDto); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusCreated
	resp.Message = "Email verification sent to your email address. Please verify your email address."

	return c.JSON(resp)
}

func (handler *authHandler) RefreshAccessToken(c *fiber.Ctx) error {
	var resp response.Response

	refreshAccessTokenRequest := new(request.RefreshAccessTokenRequest)

	if err := c.BodyParser(refreshAccessTokenRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	token, err := handler.authService.RefreshAccessToken(refreshAccessTokenRequest.RefreshToken)

	if err != nil {
		resp.Status = constants.ClientErrorBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = constants.SuccessOperationCompleted
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{"token": token}

	return c.JSON(resp)
}

func (handler *authHandler) ResendEmailVerification(c *fiber.Ctx) error {
	var resp response.Response

	emailRequest := new(request.EmailRequest)

	if err := c.BodyParser(emailRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if vEs, err := handler.validator.EmailValidate(*emailRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		resp.Data = vEs
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if err := handler.authService.ResendEmailVerification(emailRequest.Email); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = "Verification Email sent"

	return c.JSON(resp)
}

func (handler *authHandler) VerifyEmail(c *fiber.Ctx) error {
	var resp response.Response
	var req request.VerifyEmailRequest

	if err := c.BodyParser(&req); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if err := handler.authService.VerifyEmail(req.Email, req.Code); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.JSON(resp)
}

func (handler *authHandler) ForgotPassword(c *fiber.Ctx) error {
	var resp response.Response

	emailRequest := new(request.EmailRequest)

	if err := c.BodyParser(emailRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if vEs, err := handler.validator.EmailValidate(*emailRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		resp.Data = vEs
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if err := handler.authService.ForgotPassword(emailRequest.Email); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = "Reset Password Code sent to email"

	return c.JSON(resp)
}

func (handler *authHandler) ResetPassword(c *fiber.Ctx) error {
	var resp response.Response

	resetPasswordRequest := new(request.ResetPasswordRequest)

	if err := c.BodyParser(resetPasswordRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Invalid request"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if vEs, err := handler.validator.ResetPasswordValidate(*resetPasswordRequest); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		resp.Data = vEs
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if err := handler.authService.ResetPassword(resetPasswordRequest.Code, resetPasswordRequest.Email, resetPasswordRequest.Password); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.JSON(resp)
}

func (handler *authHandler) VerifyEmailAndCode(c *fiber.Ctx) error {
	var resp response.Response
	var req request.VerifyEmailAndCodeRequest

	if err := c.BodyParser(&req); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if err := handler.authService.VerifyEmailAndCode(req.Email, req.Code); err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.JSON(resp)
}
