package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/odinls/internal/errs"
)

func (h *handler) NotFound(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	return nil, errs.NewRespError(http.StatusNotFound, "")
}

func (h *handler) InternalServerError(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	return nil, errs.NewRespError(http.StatusInternalServerError, "")
}
