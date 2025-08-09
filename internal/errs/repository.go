package errs

import "errors"

var (
	ErrInvalidID             = errors.New("invalid id")
	ErrDocumentAlreadyExists = errors.New("the document already exists")
	ErrDocumentNotFound      = errors.New("can't find the document")
	ErrInvalidFloat          = errors.New("invalid float")
	ErrInvalidNumber         = errors.New("invalid number")
	ErrInvalidDate           = errors.New("invalid date")
)
