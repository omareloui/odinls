package views

import "strings"

func join(prefix, suffix string) string {
	if suffix == "" {
		return prefix
	}
	return strings.Join([]string{prefix, suffix}, "-")
}
