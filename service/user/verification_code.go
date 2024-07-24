package user_service

import (
	"time"

	"github.com/google/uuid"
	userDto "github.com/shordem/api.thryvo/dto"
	"github.com/shordem/api.thryvo/lib/helper"
	"github.com/shordem/api.thryvo/model"
	userRepositoryModule "github.com/shordem/api.thryvo/repository/user"
)

type verificationCodeService struct {
	userRepository userRepositoryModule.UserRepositoryInterface
	codeRepository userRepositoryModule.VerificationCodeRepositoryInterface
}

type VerificationCodeServiceInterface interface {
	CreateVerificationCode(email string) (string, error)
	FindCodeAndEmail(code string, email string) (userDto.VerificationCodeDTO, error)
	FindCodeByEmail(email string) (userDto.VerificationCodeDTO, error)
	DeleteVerificationCode(email string) error
	HasCodeExpired(code string) (bool, error)
}

func NewVerficationCodeService(
	userRepository userRepositoryModule.UserRepositoryInterface,
	codeRepository userRepositoryModule.VerificationCodeRepositoryInterface,
) VerificationCodeServiceInterface {
	return &verificationCodeService{
		userRepository: userRepository,
		codeRepository: codeRepository,
	}
}

func (c *verificationCodeService) ConvertToDTO(code model.VerificationCode) (codeDto userDto.VerificationCodeDTO) {

	codeDto.Code = code.Code
	codeDto.UserID = code.UserID.String()
	codeDto.CreatedAt = code.CreatedAt
	codeDto.UpdatedAt = code.UpdatedAt
	codeDto.DeletedAt = code.DeletedAt.Time

	return codeDto
}

func (c *verificationCodeService) ConvertToModel(codeDto userDto.VerificationCodeDTO) (code model.VerificationCode) {

	code.Code = codeDto.Code
	code.UserID, _ = uuid.Parse(codeDto.UserID)
	code.CreatedAt = codeDto.CreatedAt
	code.UpdatedAt = codeDto.UpdatedAt
	code.DeletedAt.Time = codeDto.DeletedAt

	return code
}

// CreateVerificationCode implements CodeServiceInterface.
func (c *verificationCodeService) CreateVerificationCode(email string) (string, error) {
	user, err := c.userRepository.FindUserByEmail(email)
	code := helper.GenerateRandomDigits(6)

	if err != nil {
		return "", err
	}

	codeModel := model.VerificationCode{
		Code:   code,
		UserID: user.ID,
	}

	newRecord, err := c.codeRepository.CreateVerificationCode(codeModel)

	if err != nil {
		return "", err
	}

	return newRecord.Code, nil
}

// DeleteVerificationCode implements CodeServiceInterface.
func (c *verificationCodeService) DeleteVerificationCode(email string) error {

	user, err := c.userRepository.FindUserByEmail(email)

	if err != nil {
		return err
	}

	return c.codeRepository.DeleteVerificationCode(user.ID)
}

// FindCodeAndEmail implements CodeServiceInterface.
func (c *verificationCodeService) FindCodeAndEmail(code string, email string) (userDto.VerificationCodeDTO, error) {

	user, err := c.userRepository.FindUserByEmail(email)

	if err != nil {
		return userDto.VerificationCodeDTO{}, err
	}

	codeModel, err := c.codeRepository.FindCodeAndUserId(code, user.ID)

	if err != nil {
		return userDto.VerificationCodeDTO{}, err
	}

	return c.ConvertToDTO(codeModel), nil

}

// FindCodeByEmail implements CodeServiceInterface.
func (c *verificationCodeService) FindCodeByEmail(email string) (userDto.VerificationCodeDTO, error) {

	user, err := c.userRepository.FindUserByEmail(email)

	if err != nil {
		return userDto.VerificationCodeDTO{}, err
	}

	codeModel, err := c.codeRepository.FindCodeByUserId(user.ID)

	if err != nil {
		return userDto.VerificationCodeDTO{}, err
	}

	return c.ConvertToDTO(codeModel), nil
}

// HasCodeExpired implements CodeServiceInterface.
func (c *verificationCodeService) HasCodeExpired(code string) (bool, error) {

	expireTimeInMinutes := 60 * time.Minute

	codeModel, err := c.codeRepository.FindByCode(code)

	if err != nil {
		return false, err
	}

	if time.Now().Before(codeModel.CreatedAt.Add(expireTimeInMinutes)) {
		return true, nil
	}

	return false, nil
}
