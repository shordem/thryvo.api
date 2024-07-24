package core_validator

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/shordem/api.thryvo/payload/request"
	"github.com/shordem/api.thryvo/validator"
)

type WishlistValidator struct {
	validator.Validator[request.CreateWishlistRequest]
}

func (validator *WishlistValidator) CreateWishlistValidate(req request.CreateWishlistRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.ProductUUID, validation.Required),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}
