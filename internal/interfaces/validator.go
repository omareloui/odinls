package interfaces

import (
	"github.com/omareloui/formmap"
)

type Validator interface {
	Validate(any) *formmap.ValidationError
}
