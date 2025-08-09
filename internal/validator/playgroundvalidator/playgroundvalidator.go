package playgroundvalidator

import (
	"regexp"
	"strings"

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

func (v *playgroundValidator) ParseError(error any) formmap.ValidationError {
	valerr := formmap.ValidationErrors{}
	for _, err := range error.(validator.ValidationErrors) {
		namespace := err.Namespace()
		firstDot := strings.Index(namespace, ".")
		path := namespace[firstDot+1:]
		valerr[path] = formmap.ValidationField{
			Tag:   err.ActualTag(),
			Param: err.Param(),
		}
	}
	return formmap.ValidationError{Errors: valerr}
}

func NewValidator() *playgroundValidator {
	val := validator.New(validator.WithRequiredStructEnabled())
	_ = val.RegisterValidation("not_blank", validators.NotBlank)
	_ = val.RegisterValidation("alphanum_with_underscore", IsAlphaNumWithUnderScore)
	return &playgroundValidator{validator: val}
}

func IsAlphaNumWithUnderScore(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	field := fl.Field()
	return re.Match([]byte(field.String()))
}
