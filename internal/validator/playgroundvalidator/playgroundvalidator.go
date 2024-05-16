package playgroundvalidator

import "github.com/go-playground/validator/v10"

type playgroundValidator struct {
	validator *validator.Validate
}

func (v *playgroundValidator) Validate(input any) error {
	return v.validator.Struct(input)
}

func NewValidator() *playgroundValidator {
	return &playgroundValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
