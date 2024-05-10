package app_errors

import "fmt"

type EntityNotFound struct {
	Identifier string
	Entity     string
}

func (e *EntityNotFound) Error() string {
	return fmt.Sprintf(`entity "%s" with filter "%s" was not found`, e.Entity, e.Identifier)
}

func NewEntityNotFound(entity string, id string) error {
	return &EntityNotFound{Entity: entity, Identifier: id}
}
