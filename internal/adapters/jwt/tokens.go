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
	Username string
	Name     string
}

func (jwt *JwtAccessClaims) String() string {
	return fmt.Sprintf("{\n  ID: \"%s\",\n  Email: \"%s\",\n  Username: \"%s\",\n  Name: \"%s\"\n}", jwt.ID, jwt.Email, jwt.Username, jwt.Name)
}

func (jwt *JwtRefreshClaims) String() string {
	return fmt.Sprintf("{ ID: \"%s\" }", jwt.ID)
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
