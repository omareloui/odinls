package handler

import (
	"context"
	"regexp"

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
