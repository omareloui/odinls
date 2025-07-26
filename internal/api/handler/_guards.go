package handler

import (
	"net/http"
)

func (h *handler) AuthGuard(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		access := r.Context().Value(authContextKey)

		if access == nil {
			if r.Header.Get("Hx-Request") == "true" {
				return hxRespondWithRedirect(w, "/")
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return nil
		}

		return next(w, r)
	}
}

func (h *handler) AlreadyAuthedGuard(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		access := r.Context().Value(authContextKey)
		if access != nil {
			if r.Header.Get("Hx-Request") == "true" {
				return hxRespondWithRedirect(w, "/")
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return nil
		}
		return next(w, r)
	}
}
