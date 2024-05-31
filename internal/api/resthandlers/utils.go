package resthandlers

import (
	"errors"
	"net/http"

	"github.com/a-h/templ"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

// TODO(refactor): update the expected signature for the handlers to return the
// error and handle the errors in an centralized place.
// Research: how to pass the validation handler
// Return all the errors from the helpers after the refactor

type ValidationErrorResponseFunc func(valerr *errs.ValidationError) templ.Component

func respondWithTemplate(w http.ResponseWriter, r *http.Request, status int, template templ.Component) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	if err := renderToBody(w, r, template); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func respondWithString(w http.ResponseWriter, r *http.Request, status int, str string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	if _, err := w.Write([]byte(str)); err != nil {
		respondWithInternalServerError(w, r)
	}
}

func respondWithInternalServerError(w http.ResponseWriter, r *http.Request) {
	respondWithErrorPage(w, r, http.StatusInternalServerError)
}

func respondWithUnauthorized(w http.ResponseWriter, r *http.Request) {
	respondWithErrorPage(w, r, http.StatusUnauthorized)
}

func respondWithForbidden(w http.ResponseWriter, r *http.Request) {
	respondWithErrorPage(w, r, http.StatusForbidden)
}

func respondWithNotFound(w http.ResponseWriter, r *http.Request) {
	respondWithErrorPage(w, r, http.StatusNotFound)
}

func respondWithErrorPage(w http.ResponseWriter, r *http.Request, status int) {
	auth := r.Context().Value(authContextKey)
	respondWithTemplate(w, r, status, views.ErrorPage(auth.(*jwtadapter.JwtAccessClaims), http.StatusText(status), status))
}

func renderToBody(w http.ResponseWriter, r *http.Request, template templ.Component) error {
	return template.Render(r.Context(), w)
}

func hxRespondWithRedirect(w http.ResponseWriter, path string) {
	w.Header().Set("HX-Location", path)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(errs.ErrForbidden, err) {
		respondWithForbidden(w, r)
		return
	}

	if errors.Is(errs.ErrInvalidID, err) || errors.Is(errs.ErrInvalidFloat, err) {
		respondWithString(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondWithInternalServerError(w, r)
}

func ErrorHandlerAdapter(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		handleError(w, r, err)
	}
}
