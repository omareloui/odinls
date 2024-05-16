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
	case "gte":
		return fmt.Sprintf("Value is low (at least %s is required)", v.Param)
	case "min":
		return fmt.Sprintf("Value is too short (at least %s characters)", v.Param)
	}
	return ""
}

type ValidationError struct {
	Errors Errors
}

func (v ValidationError) Error() string {
	return "Validation Error"
}
