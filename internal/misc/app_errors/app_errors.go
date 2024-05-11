package app_errors

import "fmt"

type EntityNotFound struct {
	Identifier string
	Entity     string
}

func (e *EntityNotFound) Error() string {
	return fmt.Sprintf(`entity "%s" with filter "%s" was not found`, e.Entity, e.Identifier)
}

func NewEntityNotFound(entity string, id string) error {
	return &EntityNotFound{Entity: entity, Identifier: id}
}

// ------------------------------------------ //

type ConfirmPasswordNotMatching struct{}

func (e *ConfirmPasswordNotMatching) Error() string {
	return "the passwords doesn't match"
}

// ------------------------------------------ //

type EmailAlreadyInUse struct{}

func (e *EmailAlreadyInUse) Error() string {
	return "email already in use"
}

// ------------------------------------------ //

type ValidationErr struct {
	Code    int16
	Message string
}

func (e *ValidationErr) Error() string {
	return fmt.Sprintf("ValidationError: %s\n", e.Message)
}

func NewValidationErr(msg string) *ValidationErr {
	return &ValidationErr{Code: 422, Message: msg}
}

// ------------------------------------------ //

type HttpErr struct {
	Code    int16
	Message string
}

func (e *HttpErr) Error() string {
	return fmt.Sprintf("HttpError: %s\n", e.Message)
}

func NewHttpErr(msg string, code int16) *HttpErr {
	return &HttpErr{
		Code:    code,
		Message: msg,
	}
}
