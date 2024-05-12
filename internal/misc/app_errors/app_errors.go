package app_errors

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type EntityNotFound struct {
	Identifier string
	Entity     string
}

func (e EntityNotFound) Error() string {
	return fmt.Sprintf(`entity "%s" with filter "%s" was not found`, e.Entity, e.Identifier)
}

func NewEntityNotFound(entity string, id string) error {
	return &EntityNotFound{Entity: entity, Identifier: id}
}

// ------------------------------------------ //

type ConfirmPasswordNotMatching struct{}

func (e ConfirmPasswordNotMatching) Error() string {
	return "the passwords doesn't match"
}

// ------------------------------------------ //

type EmailAlreadyInUse struct {
	Code   int
	Errors validator.ValidationErrors
}

func (e EmailAlreadyInUse) Error() string {
	return "email already in use"
}

func NewEmailAlreadyInUse() *EmailAlreadyInUse {
	return &EmailAlreadyInUse{Code: http.StatusUnprocessableEntity}
}

// ------------------------------------------ //

type ValidationErr struct {
	Code   int
	Errors validator.ValidationErrors
}

func (e ValidationErr) Error() string {
	return fmt.Sprintf("ValidationError: %+v\n", e.Errors)
}

func NewValidationErr(errors *validator.ValidationErrors) *ValidationErr {
	return &ValidationErr{Code: http.StatusUnprocessableEntity, Errors: *errors}
}

// ------------------------------------------ //

type HttpErr struct {
	Code    int
	Message string
}

func (e HttpErr) Error() string {
	return fmt.Sprintf("HttpError: %s\n", e.Message)
}

func NewHttpErr(msg string, code int) *HttpErr {
	return &HttpErr{
		Code:    code,
		Message: msg,
	}
}
