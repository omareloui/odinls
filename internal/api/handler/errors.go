package handler

import (
	"net/http"

	"github.com/a-h/templ"
)

func (h *handler) NotFound(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	// TODO: create  component
	return NotFound()
}

func (h *handler) InternalServerError(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	// TODO: create  component
	return InternalServerError()
}

func (h *handler) Unauthorized(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	// TODO: create  component
	return Unauthorized()
}
