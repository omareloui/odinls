package jwtadapter

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omareloui/odinls/config"
	"github.com/omareloui/odinls/internal/application/core/user"
)

type TokenType int

const (
	accessToken TokenType = iota
	refreshToken

	refreshExpiration = time.Hour * 24 * 30
	accessExpiration  = time.Minute * 15
)

type tokenDetails struct {
	Encoded    string
	Expiration time.Time
}

type TokenPair struct {
	Refresh tokenDetails
	Access  tokenDetails
}

type ClaimsPair struct {
	Refresh RefreshClaims
	Access  AccessClaims
}

type RefreshClaims struct {
	ID string
}

type AccessClaims struct {
	ID            string
	OAuthID       string
	OAuthProvider user.OAuthProvider
	Email         string
	Name          user.Name
	Picture       string
	Role          user.RoleEnum
	Craftsman     *user.Craftsman
}

func (a AccessClaims) IsCraftsman() bool {
	return a.Craftsman != nil
}

var (
	ErrInvalidTokenMethod = errors.New("invalid token method")
	ErrInvalidClaimsType  = errors.New("invalid claims type")
	ErrExpiredToken       = errors.New("expired token")
)

func NewPair(usr *user.User) (*TokenPair, error) {
	refresh, refreshExp, err := newToken(refreshToken, usr)
	if err != nil {
		return nil, err
	}
	access, accessExp, err := newToken(accessToken, usr)
	if err != nil {
		return nil, err
	}
	return &TokenPair{Access: tokenDetails{access, accessExp}, Refresh: tokenDetails{refresh, refreshExp}}, nil
}

func ParseAccessClaims(token string) (*AccessClaims, error) {
	claims, err := parse(token)
	if err != nil {
		return nil, err
	}
	_claims := *claims

	aclaims := &AccessClaims{
		ID:            _claims["id"].(string),
		OAuthID:       _claims["oauth_id"].(string),
		OAuthProvider: user.OAuthProvider(_claims["oauth_provider"].(string)),
		Role:          user.RoleEnum(_claims["role"].(float64)),
		Name: user.Name{
			First: _claims["first_name"].(string),
			Last:  _claims["last_name"].(string),
		},
		Email:   _claims["email"].(string),
		Picture: _claims["picture"].(string),
	}

	if _claims["hourly_rate"] != nil {
		aclaims.Craftsman = &user.Craftsman{
			HourlyRate: _claims["hourly_rate"].(float64),
		}
	}

	return aclaims, nil
}

func ParseRefreshClaims(token string) (*RefreshClaims, error) {
	claims, err := parse(token)
	if err != nil {
		return nil, err
	}
	_claims := *claims
	return &RefreshClaims{
		ID: _claims["id"].(string),
	}, nil
}

func parse(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidTokenMethod
		}
		return config.GetJwtSecret(), nil
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

func newToken(kind TokenType, usr *user.User) (string, time.Time, error) {
	var exp time.Time
	var claims jwt.MapClaims

	now := time.Now()

	if kind == refreshToken {
		exp = time.Now().Add(refreshExpiration)
		claims = jwt.MapClaims{
			"id": usr.ID,

			"oauth_id":       usr.OAuthID,
			"oauth_provider": usr.OAuthProvider,

			"iat": now.Unix(),
			"exp": exp.Unix(),
		}
	} else {
		exp = time.Now().Add(accessExpiration)
		claims = jwt.MapClaims{
			"id": usr.ID,

			"oauth_id":       usr.OAuthID,
			"oauth_provider": usr.OAuthProvider,

			"role": usr.Role,

			"email":      usr.Email,
			"first_name": usr.Name.First,
			"last_name":  usr.Name.Last,
			"picture":    usr.Picture,

			"iat": now.Unix(),
			"exp": exp.Unix(),
		}

		if usr.Craftsman != nil {
			claims["hourly_rate"] = usr.Craftsman.HourlyRate
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	tokenStr, err := token.SignedString(config.GetJwtSecret())
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenStr, exp, nil
}
