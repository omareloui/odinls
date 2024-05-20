package resthandlers

import (
	"context"
	"errors"
	"net/http"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
)

type authContextKeyType string

const authContextKey authContextKeyType = "auth"

var ErrNoAccessToken = errors.New("no access token context")

func (h *handler) AttachAuthenticatedUserMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(jwtadapter.AccessTokenCookieName)
		if err != nil || cookie.Value == "" {
			next.ServeHTTP(w, r)
			return
		}

		access, err := h.jwtAdapter.ParseAccessClaims(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), authContextKey, access)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func (h *handler) getAuthFromContext(r *http.Request) (*jwtadapter.JwtAccessClaims, error) {
	auth := r.Context().Value(authContextKey)
	if auth == nil {
		return nil, ErrNoAccessToken
	}
	return auth.(*jwtadapter.JwtAccessClaims), nil
}
