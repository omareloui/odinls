package resthandlers

import (
	"context"
	"errors"
	"net/http"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
)

type jwtContextKeyType string

const (
	accessTokenCookieName  = "access_token"
	refreshTokenCookieName = "refresh_token"
)

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
	accessCookie, err := r.Cookie(accessTokenCookieName)
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
	cookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil || cookie.Value == "" {
		return nil
	}

	parsed, err := h.jwtAdapter.ParseRefreshClaims(cookie.Value)
	if err != nil {
		return nil
	}

	usr, err := h.app.UserService.GetUserByID(parsed.ID)
	if err != nil {
		return nil
	}

	cookiesPair, err := h.newCookiesPairFromUser(usr)
	if err != nil {
		return nil
	}

	http.SetCookie(w, cookiesPair.Refresh)
	http.SetCookie(w, cookiesPair.Access)

	access, err := h.jwtAdapter.ParseAccessClaims(cookiesPair.Access.Value)
	if err != nil {
		return r.Context()
	}

	return context.WithValue(r.Context(), authContextKey, access)
}

func hasRefreshToken(r *http.Request) bool {
	_, err := r.Cookie(refreshTokenCookieName)
	return err == nil
}

func (h *handler) getAuthFromContext(r *http.Request) (*jwtadapter.JwtAccessClaims, error) {
	auth := r.Context().Value(authContextKey)
	if auth == nil {
		return nil, ErrNoAccessCookie
	}
	return auth.(*jwtadapter.JwtAccessClaims), nil
}
