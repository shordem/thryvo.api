package core_validator

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/shordem/api.thryvo/payload/request"
	"github.com/shordem/api.thryvo/validator"
)

type ProductValidator struct {
	validator.Validator[request.CreateProductRequest]
}

func (validator *ProductValidator) CreateProductValidate(req request.CreateProductRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.CategoryUUID, validation.Required),
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50)),
		validation.Field(&req.Description, validation.Required, validation.Length(3, 100)),
		validation.Field(&req.Specification, validation.Required, validation.Length(3, 100)),
		validation.Field(&req.Price, validation.Required, validation.Min(0)),
		validation.Field(&req.Stock, validation.Required, validation.Min(0)),
		validation.Field(&req.SlashPrice, validation.Max(req.Price)),
		validation.Field(&req.Images, validation.Required, validation.Each(validation.Required, validation.Length(3, 100))),
		validation.Field(&req.Discounts, validation.NilOrNotEmpty),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}
