package interfaces

type Validator interface {
	Validate(any) error
}
