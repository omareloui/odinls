package jwtadapter

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omareloui/odinls/internal/application/core/merchant"
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

func (jwt *JwtAccessClaims) HourlyRate() float64 {
	if jwt.CraftsmanInfo.HourlyRate > 0 {
		return jwt.CraftsmanInfo.HourlyRate
	}
	return jwt.CraftsmanInfo.Merchant.HourlyRate
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

	if r, ok := (*claims)["role"].(map[string]interface{}); ok {
		createdAt, _ := time.Parse(time.RFC3339, r["created_at"].(string))
		updatedAt, _ := time.Parse(time.RFC3339, r["updated_at"].(string))
		c.Role = role.Role{
			ID:        r["id"].(string),
			Name:      r["name"].(string),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}
	if ci, ok := (*claims)["craftsmanInfo"].(map[string]interface{}); ok {
		c.CraftsmanInfo = user.Craftsman{
			MerchantID: ci["merchant_id"].(string),
			HourlyRate: ci["hourly_rate"].(float64),
		}
		if m, ok := ci["merchant"].(map[string]interface{}); ok {
			createdAt, _ := time.Parse(time.RFC3339, m["created_at"].(string))
			updatedAt, _ := time.Parse(time.RFC3339, m["updated_at"].(string))
			c.CraftsmanInfo.Merchant = &merchant.Merchant{
				ID:        m["id"].(string),
				Name:      m["name"].(string),
				Logo:      m["logo"].(string),
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}
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
		"id":       usr.ID,
		"email":    usr.Email,
		"username": usr.Username,
		"name":     fmt.Sprintf("%s %s", usr.Name.First, usr.Name.Last),
		"exp":      exp.Unix(),
	}

	if usr.Craftsman != nil {
		claims["craftsmanInfo"] = *usr.Craftsman
	}

	if usr.Role != nil {
		claims["role"] = *usr.Role
	}

	return &claims
}

func newTokenPair(access, refresh string, accessExp, refreshExp time.Time) *TokenPair {
	return &TokenPair{
		Access:  tokenDetails{access, accessExp},
		Refresh: tokenDetails{refresh, refreshExp},
	}
}
