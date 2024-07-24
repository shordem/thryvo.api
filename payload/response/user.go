package response

import (
	"github.com/shordem/api.thryvo/dto"
)

type LoginResponse struct {
	Response

	Data dto.LoginResponseDTO `json:"data"`
}

type UserResponse struct {
	Response
	Data UserResponseData `json:"data"`
}

type UserResponseData struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	ReferralCode string `json:"referral_code"`
}
