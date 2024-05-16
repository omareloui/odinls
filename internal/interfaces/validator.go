package interfaces

import "github.com/omareloui/odinls/internal/errs"

type Validator interface {
	Validate(any) error
	ParseError(any) errs.ValidationError
}
