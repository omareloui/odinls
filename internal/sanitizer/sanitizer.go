package sanitizer

import (
	"strings"
)

func TrimString(value string) string {
	return strings.Trim(strings.TrimSpace(value), "\x1c\x1d\x1e\x1f")
}

func LowerCase(value string) string {
	return strings.ToLower(value)
}

func TrimAndLowerCaseString(value string) string {
	return TrimString(LowerCase(value))
}
