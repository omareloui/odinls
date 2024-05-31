package resthandlers

import (
	"net/http"

	"github.com/a-h/templ"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

type ValidationErrorResponseFunc func(valerr *errs.ValidationError) templ.Component

func respondWithTemplate(w http.ResponseWriter, r *http.Request, status int, template templ.Component) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	return renderToBody(w, r, template)
}

func respondWithString(w http.ResponseWriter, r *http.Request, status int, str string) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	if _, err := w.Write([]byte(str)); err != nil {
		return respondWithInternalServerError(w, r)
	}
	return nil
}

func respondWithInternalServerError(w http.ResponseWriter, r *http.Request) error {
	return respondWithErrorPage(w, r, http.StatusInternalServerError)
}

func respondWithUnauthorized(w http.ResponseWriter, r *http.Request) error {
	return respondWithErrorPage(w, r, http.StatusUnauthorized)
}

func respondWithForbidden(w http.ResponseWriter, r *http.Request) error {
	return respondWithErrorPage(w, r, http.StatusForbidden)
}

func respondWithNotFound(w http.ResponseWriter, r *http.Request) error {
	return respondWithErrorPage(w, r, http.StatusNotFound)
}

func respondWithErrorPage(w http.ResponseWriter, r *http.Request, status int) error {
	auth := r.Context().Value(authContextKey)
	return respondWithTemplate(w, r, status, views.ErrorPage(auth.(*jwtadapter.JwtAccessClaims), http.StatusText(status), status))
}

func renderToBody(w http.ResponseWriter, r *http.Request, template templ.Component) error {
	return template.Render(r.Context(), w)
}

func hxRespondWithRedirect(w http.ResponseWriter, path string) error {
	w.Header().Set("HX-Location", path)
	return nil
}
