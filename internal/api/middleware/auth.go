package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/logger"
	"go.uber.org/zap"
)

const (
	AccessClaimsCookieName  = "access_claims"
	RefreshClaimsCookieName = "refresh_claims"
)

type (
	AccessClaimsCtxKey  struct{}
	RefreshClaimsCtxKey struct{}
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var claims *jwtadapter.AccessClaims

		accessCookie, err := r.Cookie(AccessClaimsCookieName)
		if !errors.Is(err, http.ErrNoCookie) {
			accessClaims, err := jwtadapter.ParseAccessClaims(accessCookie.Value)
			if err == nil {
				claims = accessClaims
			}
		}

		if claims != nil {
			ctx := context.WithValue(r.Context(), AccessClaimsCtxKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		refreshCookie, err := r.Cookie(RefreshClaimsCookieName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		currPath := r.URL.Path
		if currPath != "/refresh-tokens" {
			queries := r.URL.Query()
			next := queries.Get("next")
			if next == "" {
				next = currPath
			}
			path := fmt.Sprintf("/refresh-tokens?next=%s", next)
			http.Redirect(w, r, path, http.StatusTemporaryRedirect)
			return
		} else {
			refreshClaims, err := jwtadapter.ParseRefreshClaims(refreshCookie.Value)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), RefreshClaimsCtxKey{}, refreshClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	})
}

func Protected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if accessClaims := r.Context().Value(AccessClaimsCtxKey{}); accessClaims == nil {
			l := logger.FromCtx(r.Context())

			httperr := errs.NewRespError(http.StatusUnauthorized, "")
			w.Header().Set("Content-Type", "text/html")
			l.Error("trying to access protected route", zap.String("path", r.URL.Path))
			http.Error(w, httperr.Message, httperr.Code)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AlreadyAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if accessClaims := r.Context().Value(AccessClaimsCtxKey{}); accessClaims != nil {
			l := logger.FromCtx(r.Context())
			l.Error("trying to access register or login page while logged", zap.String("path", r.URL.Path))

			if r.Header.Get("Hx-Request") == "true" {
				w.Header().Set("HX-Location", "/")
				return
			}

			httperr := errs.NewRespError(http.StatusUnauthorized, "")
			w.Header().Set("Content-Type", "text/html")
			http.Error(w, httperr.Message, httperr.Code)
			return
		}

		next.ServeHTTP(w, r)
	})
}
