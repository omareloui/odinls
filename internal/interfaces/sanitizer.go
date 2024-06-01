package interfaces

type Sanitizer interface {
	SanitizeStruct(interface{}) error
	Trim(str string) string
	Lower(str string) string
}
