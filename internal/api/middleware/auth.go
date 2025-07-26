package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

		currPath := r.URL.Path
		isRefreshTokenRoute := strings.Contains(currPath, "oauth/refresh-tokens")
		refreshCookie, err := r.Cookie(RefreshClaimsCookieName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if !isRefreshTokenRoute {
			queries := r.URL.Query()
			next := queries.Get("next")
			if next == "" {
				next = currPath
			}
			path := fmt.Sprintf("/oauth/refresh-tokens?next=%s", next)
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
