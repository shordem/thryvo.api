package core_validator

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/shordem/api.thryvo/payload/request"
	"github.com/shordem/api.thryvo/validator"
)

type CategoryValidator struct {
	validator.Validator[request.CreateCategoryRequest]
}

func (validator *CategoryValidator) CreateCategoryValidate(req request.CreateCategoryRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50)),
		validation.Field(&req.Description, validation.Required, validation.Length(3, 100)),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}

func (validator *CategoryValidator) UpdateCategoryValidate(req request.UpdateCategoryRequest) (map[string]interface{}, error) {
	return validator.CreateCategoryValidate(req.CreateCategoryRequest)
}
