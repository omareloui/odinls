package resthandlers

import (
	"context"
	"net/http"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
)

type authContextKeyType string

const authContextKey authContextKeyType = "auth"

func (h *handler) AttachAuthenticatedUserMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(jwtadapter.AccessTokenCookieName)
		if err != nil {
			handler(w, r)
			return
		}
		access, err := h.jwtAdapter.ParseAccessClaims(cookie.Value)
		if err != nil {
			handler(w, r)
			return
		}
		ctx := context.WithValue(context.Background(), authContextKey, access)
		handler(w, r.WithContext(ctx))
	}
}

func (h *handler) getAuthFromContext(r *http.Request) *jwtadapter.JwtAccessClaims {
	auth := r.Context().Value(authContextKey)
	return auth.(*jwtadapter.JwtAccessClaims)
}
