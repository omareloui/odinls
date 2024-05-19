package errs

import (
	"errors"
	"fmt"
)

var ErrInvalidID = errors.New("ID Invalid")

type Errors map[string]ValidationField

func (e *Errors) MsgFor(fieldName string) string {
	f, ok := (*e)[fieldName]
	if !ok {
		return ""
	}
	return f.Msg()
}

type ValidationField struct {
	Tag   string
	Param string
}

func (v ValidationField) Msg() string {
	switch v.Tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "http_url":
		return "This field must be a valid URL"
	case "gte":
		return fmt.Sprintf("Value is low (at least %s is required)", v.Param)
	case "min":
		return fmt.Sprintf("Value is too short (at least %s characters)", v.Param)
	case "max":
		return fmt.Sprintf("Value is too log (maximum %s characters)", v.Param)
	case "eqfield":
		return fmt.Sprintf(`This field must match the "%s" field`, v.Tag)
	case "not_blank":
		return "This field can't empty"
	case "alphanum":
		return "This field expects alphanumeric value"
	case "alphanum_with_underscore":
		return "This field expects alphanumeric or underscore characters"
	case "mongodb":
		return "Invalid ID"
	default:
		msg := fmt.Sprintf(`Failed on "%s" tag`, v.Tag)
		if v.Param != "" {
			msg += fmt.Sprintf(` with "%s" as a param`, v.Param)
		}
		return msg
	}
}

type ValidationError struct {
	Errors Errors
}

func (v ValidationError) Error() string {
	return "Validation Error"
}
