package conformadaptor

import (
	"strings"

	"github.com/leebenson/conform"
)

type sanitizer struct{}

func NewSanitizer() *sanitizer {
	return &sanitizer{}
}

func (s *sanitizer) SanitizeStruct(in interface{}) error {
	return conform.Strings(in)
}

func (s *sanitizer) Trim(str string) string {
	return strings.Trim(strings.TrimSpace(str), "\x1c\x1d\x1e\x1f")
}

func (s *sanitizer) Lower(str string) string {
	return strings.ToLower(str)
}

func (s *sanitizer) Upper(str string) string {
	return strings.ToUpper(str)
}
