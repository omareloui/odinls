package resthandlers

import (
	"context"
	"errors"
	"net/http"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
)

type jwtContextKeyType string

const (
	authContextKey    jwtContextKeyType = "auth"
	refreshContextKey jwtContextKeyType = "refresh"
)

var ErrNoAccessCookie = errors.New("no access token context")

func (h *handler) AttachAuthenticatedUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(h.getAuthedContext(w, r)))
	})
}

func (h *handler) getAuthedContext(w http.ResponseWriter, r *http.Request) context.Context {
	accessCookie, err := r.Cookie(jwtadapter.AccessTokenCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) && hasRefreshToken(r) {
			return h.refreshTokensAndGetContext(w, r)
		}
		return r.Context()
	}

	parsed, err := h.jwtAdapter.ParseAccessClaims(accessCookie.Value)
	if err != nil {
		if errors.Is(err, jwtadapter.ErrExpiredToken) && hasRefreshToken(r) {
			return h.refreshTokensAndGetContext(w, r)
		}
		return r.Context()
	}

	return context.WithValue(r.Context(), authContextKey, parsed)
}

func (h *handler) refreshTokensAndGetContext(w http.ResponseWriter, r *http.Request) context.Context {
	cookies := h.refreshTokens(w, r)

	access, err := h.jwtAdapter.ParseAccessClaims(cookies.Access.Value)
	if err != nil {
		return r.Context()
	}

	return context.WithValue(r.Context(), authContextKey, access)
}

func hasRefreshToken(r *http.Request) bool {
	_, err := r.Cookie(jwtadapter.RefreshTokenCookieName)
	return err == nil
}

func (h *handler) getAuthFromContext(r *http.Request) (*jwtadapter.JwtAccessClaims, error) {
	auth := r.Context().Value(authContextKey)
	if auth == nil {
		return nil, ErrNoAccessCookie
	}
	return auth.(*jwtadapter.JwtAccessClaims), nil
}
