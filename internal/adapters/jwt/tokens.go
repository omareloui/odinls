package jwtadapter

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omareloui/odinls/internal/application/core/user"
)

type (
	TokenPair struct {
		Refresh tokenDetails
		Access  tokenDetails
	}
	tokenDetails struct {
		Encoded    string
		Expiration time.Time
	}
)

type (
	JwtClaimsPair struct {
		Refresh JwtRefreshClaims
		Access  JwtAccessClaims
	}

	JwtRefreshClaims struct {
		ID string
	}

	JwtAccessClaims struct {
		ID            string
		Email         string
		Username      string
		FirstName     string
		LastName      string
		Role          user.RoleEnum
		CraftsmanInfo *user.Craftsman
	}
)

func (jwt *JwtAccessClaims) IsCraftsman() bool {
	return jwt.CraftsmanInfo != nil
}

func (jwt *JwtAccessClaims) HourlyRate() float64 {
	return jwt.CraftsmanInfo.HourlyRate
}

func newRefreshClaimsFromMapClaims(claims *jwt.MapClaims) *JwtRefreshClaims {
	return &JwtRefreshClaims{ID: (*claims)["id"].(string)}
}

func newAccessClaimsFromMapClaims(claims *jwt.MapClaims) *JwtAccessClaims {
	c := &JwtAccessClaims{
		ID:        (*claims)["id"].(string),
		FirstName: (*claims)["first_name"].(string),
		LastName:  (*claims)["last_name"].(string),
		Email:     (*claims)["email"].(string),
		Username:  (*claims)["username"].(string),
		Role:      (*claims)["role"].(user.RoleEnum),
	}

	if ci, ok := (*claims)["craftsmanInfo"].(map[string]any); ok {
		c.CraftsmanInfo = &user.Craftsman{
			HourlyRate: ci["hourly_rate"].(float64),
		}
	}

	return c
}

func newRefreshMapClaims(usr *user.User, exp time.Time) *jwt.MapClaims {
	return &jwt.MapClaims{
		"id":  usr.ID,
		"exp": exp.Unix(),
	}
}

func newAccessMapClaims(usr *user.User, exp time.Time) *jwt.MapClaims {
	claims := jwt.MapClaims{
		"id":         usr.ID,
		"first_name": usr.Name.First,
		"last_name":  usr.Name.Last,
		"email":      usr.Email,
		"username":   usr.Username,
		"role":       usr.Role,
		"exp":        exp.Unix(),
	}

	if usr.Craftsman != nil {
		claims["craftsmanInfo"] = *usr.Craftsman
	}

	return &claims
}

func newTokenPair(access, refresh string, accessExp, refreshExp time.Time) *TokenPair {
	return &TokenPair{
		Access:  tokenDetails{access, accessExp},
		Refresh: tokenDetails{refresh, refreshExp},
	}
}
