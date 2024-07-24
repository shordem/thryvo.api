package userResponse

import "github.com/shordem/api.thryvo/payload/response"

type UserResponse struct {
	response.Response
	Data UserResponseData `json:"data"`
}

type UserResponseData struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	ReferralCode string `json:"referral_code"`
}
