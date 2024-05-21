package resthandlers

import (
	"net/http"
)

func (h *handler) AuthGuard(next http.HandlerFunc) http.HandlerFunc {
	// TODO(auth): try to refresh the token here?
	return func(w http.ResponseWriter, r *http.Request) {
		access := r.Context().Value(authContextKey)
		if access == nil {
			if r.Header.Get("Hx-Request") == "true" {
				hxRespondWithRedirect(w, "/")
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			return
		}
		next(w, r)
	}
}

func (h *handler) AlreadyAuthedGuard(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		access := r.Context().Value(authContextKey)
		if access != nil {
			if r.Header.Get("Hx-Request") == "true" {
				hxRespondWithRedirect(w, "/")
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			return
		}
		next(w, r)
	}
}
