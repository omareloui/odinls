package counter

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type CounterService interface {
	AddOneToProduct(claims *jwtadapter.AccessClaims, category string) (uint8, error)
	AddOneToOrder(claims *jwtadapter.AccessClaims) (uint, error)
}
