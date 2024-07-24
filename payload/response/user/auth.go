package userResponse

import (
	"github.com/shordem/api.thryvo/dto"
	"github.com/shordem/api.thryvo/payload/response"
)

type LoginResponse struct {
	response.Response

	Data dto.LoginResponseDTO `json:"data"`
}
