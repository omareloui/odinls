package jwtadapter

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omareloui/odinls/internal/application/core/user"
)

type TokenPair struct {
	Access  string
	Refresh string
}

type CookiePair struct {
	Access  *http.Cookie
	Refresh *http.Cookie
}

type JwtAccessClaims struct {
	ID       string
	Email    string
	Name     string
	Username string
}

func NewAccessClaims(usr *user.User) *JwtAccessClaims {
	return &JwtAccessClaims{
		ID:       usr.ID,
		Email:    usr.Email,
		Username: usr.Username,
		Name:     fmt.Sprintf("%s %s", usr.Name.First, usr.Name.Last),
	}
}

func NewAccessClaimsFromClaims(claims *jwt.MapClaims) *JwtAccessClaims {
	return &JwtAccessClaims{
		ID:       (*claims)["id"].(string),
		Name:     (*claims)["name"].(string),
		Email:    (*claims)["email"].(string),
		Username: (*claims)["username"].(string),
	}
}

type JwtRefreshClaims struct {
	ID string
}

func NewRefreshClaims(id string) *JwtRefreshClaims {
	return &JwtRefreshClaims{ID: id}
}

func NewRefreshClaimsFromClaims(claims *jwt.MapClaims) *JwtRefreshClaims {
	return &JwtRefreshClaims{ID: (*claims)["id"].(string)}
}

type JwtClaimsPair struct {
	Access  JwtAccessClaims
	Refresh JwtRefreshClaims
}

func newTokenPair(access, refresh string) *TokenPair {
	return &TokenPair{
		Access:  access,
		Refresh: refresh,
	}
}
