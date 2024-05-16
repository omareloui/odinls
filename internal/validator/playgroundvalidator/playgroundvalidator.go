package playgroundvalidator

import (
	"github.com/go-playground/validator/v10"
	"github.com/omareloui/odinls/internal/errs"
)

type playgroundValidator struct {
	validator *validator.Validate
}

func (v *playgroundValidator) Validate(input any) error {
	return v.validator.Struct(input)
}

func (v *playgroundValidator) ParseError(error any) errs.ValidationError {
	valerr := errs.Errors{}
	for _, err := range error.(validator.ValidationErrors) {
		valerr[err.Field()] = errs.ValidationField{
			Tag:   err.ActualTag(),
			Param: err.Param(),
		}
	}
	return errs.ValidationError{Errors: valerr}
}

func NewValidator() *playgroundValidator {
	return &playgroundValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
