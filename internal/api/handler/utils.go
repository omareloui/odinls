package handler

import (
	"context"
	"net/http"
	"regexp"

	"github.com/a-h/templ"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/api/middleware"
)

var nonDigitRegexp *regexp.Regexp = regexp.MustCompile(`\D+`)

func getClaims(ctx context.Context) *jwtadapter.AccessClaims {
	v := ctx.Value(middleware.AccessClaimsCtxKey{})
	if v == nil {
		return nil
	}
	return v.(*jwtadapter.AccessClaims)
}

func respondWithTemplate(w http.ResponseWriter, r *http.Request, status int, template templ.Component) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	return renderToBody(w, r, template)
}

func respondWithString(w http.ResponseWriter, r *http.Request, status int, str string) error {
	w.Header().Set("Content-Type", "text/plain")
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
	// TODO:
	// auth := r.Context().Value(middleware.AccessClaimsCtxKey)
	// return respondWithTemplate(w, r, status, views.ErrorPage(auth.(*jwtadapter.AccessClaims), http.StatusText(status), status))
	return nil
}

func renderToBody(w http.ResponseWriter, r *http.Request, template templ.Component) error {
	return template.Render(r.Context(), w)
}
