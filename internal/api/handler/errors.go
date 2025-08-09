package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/odinls/internal/api/responder"
)

func (h *handler) NotFound(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	// TODO: create component
	return responder.NotFound()
}

func (h *handler) InternalServerError(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	// TODO: create component
	return responder.InternalServerError()
}

func (h *handler) Unauthorized(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	// TODO: create component
	return responder.Unauthorized()
}
