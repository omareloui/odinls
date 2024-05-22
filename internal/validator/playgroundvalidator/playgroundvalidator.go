package playgroundvalidator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
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
	val := validator.New(validator.WithRequiredStructEnabled())
	val.RegisterValidation("not_blank", validators.NotBlank)
	val.RegisterValidation("alphanum_with_underscore", IsAlphaNumWithUnderScore)
	return &playgroundValidator{validator: val}
}

func IsAlphaNumWithUnderScore(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	field := fl.Field()
	return re.Match([]byte(field.String()))
}
