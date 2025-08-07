package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetHomepage(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	claims := getClaims(r.Context())
	comp := views.Homepage(claims)
	return RespondOK(w, RespondWithComponent(comp))
}
