package jwtadapter

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omareloui/odinls/internal/application/core/role"
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
		Name          string
		Role          role.Role
		CraftsmanInfo user.Craftsman
	}
)

func (jwt *JwtAccessClaims) IsCraftsman() bool {
	return jwt.CraftsmanInfo.MerchantID != ""
}

func newRefreshClaimsFromMapClaims(claims *jwt.MapClaims) *JwtRefreshClaims {
	return &JwtRefreshClaims{ID: (*claims)["id"].(string)}
}

func newAccessClaimsFromMapClaims(claims *jwt.MapClaims) *JwtAccessClaims {
	c := &JwtAccessClaims{
		ID:       (*claims)["id"].(string),
		Name:     (*claims)["name"].(string),
		Email:    (*claims)["email"].(string),
		Username: (*claims)["username"].(string),
	}

	if role, ok := (*claims)["role"].(role.Role); ok {
		c.Role = role
	}
	if craftsmanInfo, ok := (*claims)["craftsmanInfo"].(user.Craftsman); ok {
		c.CraftsmanInfo = craftsmanInfo
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
	return &jwt.MapClaims{
		"id":            usr.ID,
		"email":         usr.Email,
		"username":      usr.Username,
		"name":          fmt.Sprintf("%s %s", usr.Name.First, usr.Name.Last),
		"role":          usr.Role,
		"craftsmanInfo": usr.Craftsman,
		"exp":           exp.Unix(),
	}
}

func newTokenPair(access, refresh string, accessExp, refreshExp time.Time) *TokenPair {
	return &TokenPair{
		Access:  tokenDetails{access, accessExp},
		Refresh: tokenDetails{refresh, refreshExp},
	}
}
