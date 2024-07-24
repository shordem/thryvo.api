package user_service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/shordem/api.thryvo/dto"
	"github.com/shordem/api.thryvo/lib/constants"
	"github.com/shordem/api.thryvo/lib/helper"
	"github.com/shordem/api.thryvo/service"
)

var (
	ErrEmailNotVerifed = errors.New("email is not verified")
)

type authService struct {
	userService UserServiceInterface
	codeService VerificationCodeServiceInterface
	encrpyt     helper.HashingInterface
	auth        helper.AuthInterface
	mail        service.EmailServiceInterface
}

type AuthServiceInterface interface {
	CheckEmail(email string) (uint16, error)
	Login(email, password string) (dto.LoginResponseDTO, uint16, error)
	Register(authDto dto.AuthDTO) error
	RefreshAccessToken(refreshToken string) (string, error)
	ResendEmailVerification(email string) error
	VerifyEmail(email string, code string) error
	ForgotPassword(email string) error
	ResetPassword(code, email, password string) error
	VerifyEmailAndCode(email, code string) error
}

func NewAuthService(
	userService UserServiceInterface,
	codeService VerificationCodeServiceInterface,
	mailService service.EmailServiceInterface,
) AuthServiceInterface {
	return &authService{
		userService: userService,
		codeService: codeService,
		encrpyt:     helper.NewHashing(),
		auth:        helper.NewAuth(),
		mail:        mailService,
	}
}

func (service *authService) CheckEmail(email string) (uint16, error) {
	_, err := service.userService.FindUserByEmail(email)

	if err == gorm.ErrRecordNotFound {
		return constants.UserNotFound, err
	}

	if err != nil {
		return constants.ServerErrorInternal, err
	}

	return constants.SuccessOperationCompleted, nil
}

func (service *authService) Login(email, password string) (dto.LoginResponseDTO, uint16, error) {
	user, err := service.userService.FindUserByEmail(email)

	if err == gorm.ErrRecordNotFound {
		return dto.LoginResponseDTO{}, constants.UserNotFound, errors.New("user not found")
	}

	if err != nil {
		return dto.LoginResponseDTO{}, constants.ServerErrorInternal, err
	}

	match, err := service.encrpyt.ComparePassword(password, user.Password)

	if err != nil {
		return dto.LoginResponseDTO{}, constants.ServerErrorInternal, err
	}

	if !match {
		return dto.LoginResponseDTO{}, constants.InvalidCredentials, errors.New("invalid password")
	}

	if !user.IsEmailVerified {
		return dto.LoginResponseDTO{}, constants.AccountVerificationRequired, ErrEmailNotVerifed
	}

	accessToken, err := service.auth.CreateToken(user.ID.String(), "access")

	if err != nil {
		return dto.LoginResponseDTO{}, constants.ServerErrorInternal, err
	}

	refreshToken, err := service.auth.CreateToken(user.ID.String(), "refresh")

	if err != nil {
		return dto.LoginResponseDTO{}, constants.ServerErrorInternal, err
	}

	tokenDto := dto.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokenDto, constants.SuccessOperationCompleted, nil
}

func (service *authService) Register(authDto dto.AuthDTO) error {
	var userDto dto.UserDTO

	hash, err := service.encrpyt.HashPassword(authDto.Password)

	if err != nil {
		return err
	}

	userDto.FirstName = authDto.FirstName
	userDto.LastName = authDto.LastName
	userDto.Email = authDto.Email
	userDto.Password = hash
	userDto.IsEmailVerified = false
	userDto.Role = UserRoleCustomer

	newUser, err := service.userService.CreateUser(userDto)

	if err != nil {
		service.userService.DeleteUser(newUser.ID)
		return err
	}

	err = service.SendEmail(authDto.Email, "confirm-email")

	if err != nil {
		return err
	}

	return nil
}

// RefreshAccessToken implements AuthServiceInterface.
func (service *authService) RefreshAccessToken(refreshToken string) (string, error) {
	userId, err := service.auth.ExtractUserID(refreshToken, "refresh")

	if err != nil {
		return "", err
	}

	accessToken, err := service.auth.CreateToken(userId.String(), "access")

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// VerifyEmail implements AuthServiceInterface.
func (service *authService) VerifyEmail(email string, code string) error {
	user, err := service.userService.FindUserByEmail(email)

	if err != nil {
		return err
	}

	if err := service.VerifyEmailAndCode(email, code); err != nil {
		return err
	}

	if err := service.codeService.DeleteVerificationCode(email); err != nil {
		return err
	}

	if user.IsEmailVerified {
		return errors.New("email already verified")
	}

	user.IsEmailVerified = true

	if _, err := service.userService.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

// ForgotPassword implements AuthServiceInterface.
func (service *authService) ForgotPassword(email string) error {

	if err := service.SendEmail(email, "reset-password"); err != nil {
		return err
	}

	return nil
}

// ResendEmailVerification implements AuthServiceInterface.
func (service *authService) ResendEmailVerification(email string) error {

	if err := service.SendEmail(email, "confirm-email"); err != nil {
		return err
	}

	return nil
}

// ResetPassword implements AuthServiceInterface.
func (service *authService) ResetPassword(code string, email string, password string) error {

	if err := service.VerifyEmailAndCode(email, code); err != nil {
		return err
	}

	if err := service.codeService.DeleteVerificationCode(email); err != nil {
		return err
	}

	user, err := service.userService.FindUserByEmail(email)

	if err != nil {
		return err
	}

	hash, err := service.encrpyt.HashPassword(password)

	if err != nil {
		return err
	}

	user.Password = hash

	if _, err := service.userService.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *authService) SendEmail(email string, templateType string) error {
	var sendEmailParams service.SendEmailParams

	user, err := s.userService.FindUserByEmail(email)

	if err != nil {
		return err
	}

	customer, err := s.userService.FindUserById(user.ID.String())

	if err != nil {
		return err
	}

	codeExists, _ := s.codeService.FindCodeByEmail(email)

	if codeExists.Code != "" {
		if err := s.codeService.DeleteVerificationCode(user.Email); err != nil {
			return err
		}
	}

	if user.IsEmailVerified && templateType == "confirm-email" {
		return errors.New("email already verified")
	}

	code, err := s.codeService.CreateVerificationCode(email)

	if err != nil {
		return err
	}

	emailVars := map[string]interface{}{
		"FullName": customer.FirstName + " " + customer.LastName,
		"Code":     []string{code},
	}

	sendEmailParams.To = email
	sendEmailParams.Template = templateType
	sendEmailParams.Variables = emailVars
	switch templateType {
	case "confirm-email":
		sendEmailParams.Subject = "Thanks for Signing Up on Mazimart"
		_ = s.mail.SendEmail(sendEmailParams)

	case "reset-password":
		sendEmailParams.Subject = "Mazimart Reset Password Request"
		_ = s.mail.SendEmail(sendEmailParams)
	}

	return nil
}

func (service *authService) VerifyEmailAndCode(email, code string) error {
	_, err := service.codeService.FindCodeAndEmail(code, email)

	if err != nil {
		return err
	}

	codeExpired, err := service.codeService.HasCodeExpired(code)

	if err != nil {
		return err
	}

	if !codeExpired {
		return errors.New("code expired")
	}

	return nil
}
