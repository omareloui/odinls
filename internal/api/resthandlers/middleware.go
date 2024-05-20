package resthandlers

import (
	"context"
	"fmt"
	"net/http"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
)

type authContextKeyType string

const authContextKey authContextKeyType = "auth"

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

		fmt.Println(ctx.Value(authContextKey))

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func (h *handler) getAuthFromContext(r *http.Request) *jwtadapter.JwtAccessClaims {
	auth := r.Context().Value(authContextKey)
	return auth.(*jwtadapter.JwtAccessClaims)
}
