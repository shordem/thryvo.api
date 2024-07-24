package core_validator

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/shordem/api.thryvo/payload/request"
	"github.com/shordem/api.thryvo/validator"
)

type ReviewValidator struct {
	validator.Validator[request.CreateReviewRequest]
}

func (validator *ReviewValidator) CreateReviewValidate(req request.CreateReviewRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.ProductID, validation.Required),
		validation.Field(&req.Rating, validation.Required, validation.Min(1), validation.Max(5)),
		validation.Field(&req.Comment, validation.Required, validation.Length(3, 512)),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}
