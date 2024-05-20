package jwtadapter

// TODO(refactor): remove unused methods

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omareloui/odinls/internal/application/core/user"
)

type TokenType int

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

const (
	AccessToken TokenType = iota + 1
	RefreshToken
)

const (
	refreshExpirationInHours  = 24 * 30
	accessExpirationInMinutes = 15
)

var (
	ErrInvalidTokenMethod = errors.New("invalid token method")
	ErrInvalidClaimsType  = errors.New("invalid claims type")
	ErrExpiredToken       = errors.New("expired token")
)

type JwtAdapter interface {
	GenPair(usr *user.User) (*TokenPair, error)
	ParsePair(accessToken, refreshToken string) (*JwtClaimsPair, error)
	GenTokenPairInCookie(usr *user.User) (*CookiePair, error)
	ParseAccessClaims(token string) (*JwtAccessClaims, error)
	ParseRefreshClaims(token string) (*JwtRefreshClaims, error)
	// TODO(refactor): make the cookies in its own adapter?
	NewCookie(tokenType TokenType, token string, expr time.Time) *http.Cookie
}

type JwtV5Adapter struct {
	secret        []byte
	signingMethod jwt.SigningMethod
}

func NewJWTV5Adapter(secret []byte) *JwtV5Adapter {
	return &JwtV5Adapter{
		secret:        secret,
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (a *JwtV5Adapter) GenPair(usr *user.User) (*TokenPair, error) {
	access, err := a.genAccessToken(usr)
	if err != nil {
		return nil, err
	}
	refresh, err := a.genRefreshToken(usr)
	if err != nil {
		return nil, err
	}
	return newTokenPair(access, refresh), nil
}

func (a *JwtV5Adapter) ParsePair(accessToken, refreshToken string) (*JwtClaimsPair, error) {
	accessClaims, err := a.ParseAccessClaims(accessToken)
	if err != nil {
		return nil, err
	}
	refreshClaims, err := a.ParseRefreshClaims(refreshToken)
	if err != nil {
		return nil, err
	}
	return &JwtClaimsPair{Access: *accessClaims, Refresh: *refreshClaims}, nil
}

func (a *JwtV5Adapter) ParseAccessClaims(token string) (*JwtAccessClaims, error) {
	claims, err := a.parse(token)
	if err != nil {
		return nil, err
	}
	return NewAccessClaimsFromClaims(claims), nil
}

func (a *JwtV5Adapter) ParseRefreshClaims(token string) (*JwtRefreshClaims, error) {
	claims, err := a.parse(token)
	if err != nil {
		return nil, err
	}
	return NewRefreshClaimsFromClaims(claims), nil
}

func (a *JwtV5Adapter) NewCookie(tokenType TokenType, token string, expr time.Time) *http.Cookie {
	tokenName := AccessTokenCookieName

	if tokenType == RefreshToken {
		tokenName = RefreshTokenCookieName
	}

	return &http.Cookie{
		Name:     tokenName,
		Value:    token,
		HttpOnly: true,
		Expires:  expr,
		Path:     "/",
	}
}

func (a *JwtV5Adapter) GenTokenPairInCookie(usr *user.User) (*CookiePair, error) {
	accessExp := time.Now().Add(time.Minute * accessExpirationInMinutes)
	refreshExp := time.Now().Add(time.Hour * refreshExpirationInHours)

	accessClaims := jwt.MapClaims{
		"id":       usr.ID,
		"email":    usr.Email,
		"username": usr.Username,
		"name":     fmt.Sprintf("%s %s", usr.Name.First, usr.Name.Last),
		"exp":      accessExp.Unix(),
	}
	refreshClaims := jwt.MapClaims{
		"id":  usr.ID,
		"exp": refreshExp.Unix(),
	}

	accessToken, err := a.genToken(accessClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.genToken(refreshClaims)
	if err != nil {
		return nil, err
	}

	accessCookie := a.NewCookie(AccessToken, accessToken, accessExp)
	refreshCookie := a.NewCookie(RefreshToken, refreshToken, refreshExp)

	return &CookiePair{Access: accessCookie, Refresh: refreshCookie}, nil
}

func (a *JwtV5Adapter) parse(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidTokenMethod
		}
		return a.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, ErrInvalidClaimsType
	}

	return &claims, nil
}

func (a *JwtV5Adapter) genAccessToken(usr *user.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       usr.ID,
		"email":    usr.Email,
		"username": usr.Username,
		"name":     fmt.Sprintf("%s %s", usr.Name.First, usr.Name.Last),
		"exp":      time.Now().Add(time.Minute * accessExpirationInMinutes).Unix(),
	}
	return a.genToken(claims)
}

func (a *JwtV5Adapter) genRefreshToken(usr *user.User) (string, error) {
	claims := jwt.MapClaims{
		"id":  usr.ID,
		"exp": time.Now().Add(time.Hour * refreshExpirationInHours).Unix(),
	}
	return a.genToken(claims)
}

func (a *JwtV5Adapter) genToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(a.signingMethod, claims)
	return token.SignedString(a.secret)
}
