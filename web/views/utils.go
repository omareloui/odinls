package views

import (
	"encoding/json"
	"strings"
)

func join(prefix, suffix string) string {
	if suffix == "" {
		return prefix
	}
	return strings.Join([]string{prefix, suffix}, "-")
}

func toJSON(in any) string {
	json, _ := json.Marshal(in)
	return string(json)
}
