package jwtadapter

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omareloui/odinls/internal/application/core/user"
)

type TokenType int

const (
	accessToken TokenType = iota
	refreshToken
)

const (
	refreshExpiration = time.Hour * 24 * 30
	accessExpiration  = time.Minute * 15
)

var (
	ErrInvalidTokenMethod = errors.New("invalid token method")
	ErrInvalidClaimsType  = errors.New("invalid claims type")
	ErrExpiredToken       = errors.New("expired token")
)

type JwtAdapter interface {
	NewPair(usr *user.User) (*TokenPair, error)
	ParseAccessClaims(token string) (*JwtAccessClaims, error)
	ParseRefreshClaims(token string) (*JwtRefreshClaims, error)
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

func (a *JwtV5Adapter) NewPair(usr *user.User) (*TokenPair, error) {
	refresh, refreshExp, err := a.newToken(refreshToken, usr)
	if err != nil {
		return nil, err
	}
	access, accessExp, err := a.newToken(accessToken, usr)
	if err != nil {
		return nil, err
	}
	return newTokenPair(access, refresh, accessExp, refreshExp), nil
}

func (a *JwtV5Adapter) ParseAccessClaims(token string) (*JwtAccessClaims, error) {
	claims, err := a.parse(token)
	if err != nil {
		return nil, err
	}
	return newAccessClaimsFromMapClaims(claims), nil
}

func (a *JwtV5Adapter) ParseRefreshClaims(token string) (*JwtRefreshClaims, error) {
	claims, err := a.parse(token)
	if err != nil {
		return nil, err
	}
	return newRefreshClaimsFromMapClaims(claims), nil
}

func (a *JwtV5Adapter) parse(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidTokenMethod
		}
		return a.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaimsType
	}

	return &claims, nil
}

func (a *JwtV5Adapter) newToken(kind TokenType, usr *user.User) (string, time.Time, error) {
	var exp time.Time
	var claims *jwt.MapClaims

	if kind == refreshToken {
		exp = time.Now().Add(refreshExpiration)
		claims = newRefreshMapClaims(usr, exp)
	} else {
		exp = time.Now().Add(accessExpiration)
		claims = newAccessMapClaims(usr, exp)
	}

	token := jwt.NewWithClaims(a.signingMethod, claims)

	tokenStr, err := token.SignedString(a.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenStr, exp, nil
}
