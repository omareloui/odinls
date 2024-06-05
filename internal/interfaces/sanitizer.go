package interfaces

type Sanitizer interface {
	SanitizeStruct(interface{}) error
	Trim(str string) string
	TrimMap(m map[string]string) map[string]string
	Lower(str string) string
	Upper(str string) string
}
